package crapchat

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"appengine"
	"appengine/user"
)

func Serve() {
	r := mux.NewRouter()

	r.HandleFunc("/logout", LogoutHandler)
	r.HandleFunc("/api/me", MeHandler).Methods("GET")
	r.HandleFunc("/api/me/friend/{email}", AddFriendHandler).Methods("POST")
	r.HandleFunc("/api/snap", AddFriendHandler).Methods("POST")

	// Static Server
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	n := negroni.New(
		negroni.HandlerFunc(AuthMiddleware),
	)
	n.UseHandler(r)

	http.Handle("/", n)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	url, _ := user.LogoutURL(c, "/")
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
	return
}

// LoginRequiredMiddleware ensures a User is logged in, otherwise redirects them to the login page.
func AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	log.Printf("%+v\n", u)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}

	next(w, r)
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	me := GetOrCreateUser(u.String())

	w.Header().Set("Content-Type", "application/json")
	w.Write(me.ToJSON())
}

func AddFriendHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	me := GetOrCreateUser(u.String())

	vars := mux.Vars(r)
	id, ok := vars["email"]
	if !ok {
		http.Error(w, "Error adding friend", 400)
		return
	}

	me.AddFriend(id)

	w.Header().Set("Content-Type", "application/json")
	w.Write(me.ToJSON())
}

func SendSnapHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	me := GetOrCreateUser(u.String())

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Body must contain to and media field", 400)
		return
	}

	crap, err := CrapFromJSON(body)
	if err != nil {
		http.Error(w, "Not valid json", 400)
	}

	if crap.To == "" {
		http.Error(w, "Not a valid 'to' field", 400)
		return
	} else if crap.Media == "" {
		http.Error(w, "Not a valid 'media' field", 400)
		return
	}

	tos := strings.Split(crap.To, ",")
	media := crap.Media
	SendCrap(tos, me.Username, media)

	fmt.Fprintf(w, "Snap sent!")
}
