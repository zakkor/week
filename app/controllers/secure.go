package controllers

import (
	"github.com/revel/revel"
	"math/rand"
	"strconv"
	"time"
)

type Secure struct {
	*revel.Controller
}

func (c *Secure) checkSession() revel.Result {
	if localsid, ok := userSessions[c.Session["user"]]; ok {
		if localsid == c.Session["sid"] {
			return nil
		}
	}
	revel.INFO.Println("did not have valid session!!!!")

	c.Flash.Error("You must log in or register")
	return c.Redirect(User.SignIn)
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

func (c *Secure) Exclusive() revel.Result {
	if res := c.checkSession(); res != nil {
		return res
	}
	return c.Render()
}
