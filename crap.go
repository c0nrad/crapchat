package crapchat

import (
	"encoding/json"
	"time"
)

type Crap struct {
	TS time.Time

	Media string

	To   string
	From string

	Seen bool
}

func CrapFromJSON(in []byte) (Crap, error) {
	var crap Crap
	err := json.Unmarshal(in, &crap)
	return crap, err
}

var Craps []*Crap

func newCrap(to, from, media string) *Crap {
	crap := &Crap{TS: time.Now(), Media: media, To: to, From: from, Seen: false}
	Craps = append(Craps, crap)
	return crap
}

func SendCrap(tos []string, from, media string) {
	for _, to := range tos {
		newCrap(to, from, media)
	}
}

func GetCraps(to string) []*Crap {
	out := []*Crap{}
	for _, crap := range Craps {
		if crap.To == to {
			out = append(out, crap)
		}
	}
	return out
}
