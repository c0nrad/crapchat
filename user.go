package crapchat

import "encoding/json"

type User struct {
	Username string
	Friends  []string
}

var Users []*User

func newUser(username string) *User {
	user := &User{username, []string{}}
	Users = append(Users, user)
	return user
}

func GetOrCreateUser(username string) *User {
	for _, user := range Users {
		if user.Username == username {
			return user
		}
	}
	return newUser(username)
}

func (user *User) AddFriend(friend string) {
	user.Friends = append(user.Friends, friend)
}

func (user *User) ToJSON() []byte {
	dater, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	return dater
}
