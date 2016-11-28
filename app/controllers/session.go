package controllers

import (
	"github.com/revel/revel"
	"math/rand"
	"strconv"
	"time"
)

// TODO: save and load sessions to db
var userSessions map[string]string = make(map[string]string)

func checkSession(c *revel.Session, f *revel.Flash) bool {
	if localsid, ok := userSessions[(*c)["user"]]; ok {
		if localsid == (*c)["sid"] {
			return true
		}
	}

	revel.INFO.Println("invalid session!!!!")
	// TODO: report to the fbi

	(*f).Error("You must log in or register")
	return false
}

func generateSession(c *revel.Session) {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	randFloat := r.Float64()
	randString := strconv.FormatFloat(randFloat, 'E', -1, 64)

	newSid := hash(randString)
	userSessions[(*c)["user"]] = newSid
	(*c)["sid"] = newSid
}
