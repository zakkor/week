package controllers

import (
	"github.com/revel/revel"
	"math/rand"
	"strconv"
	"time"
	"github.com/zakkor/week/app"
)

func checkSession(c *revel.Session, f *revel.Flash) bool {
	res, err := app.DB.Queryx("SELECT session FROM users WHERE name=$1",
		(*c)["user"])
	if err != nil {
		revel.INFO.Println(err)
		return false
	}
	res.Next()

	var dbSid string
	res.Scan(&dbSid)

	if dbSid == (*c)["sid"] && dbSid != "" && (*c)["sid"] != "" {
		return true
	}

	revel.INFO.Println("invalid session!!!!")
	// TODO: report to the fbi

	f.Error("You must log in or register")
	return false
}

func generateSession(c *revel.Session) {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	randFloat := r.Float64()
	randString := strconv.FormatFloat(randFloat, 'E', -1, 64)

	newSid := hash(randString)

	//insert into db
	_, err := app.DB.Exec("UPDATE users SET session = $1 WHERE name = $2;",
		newSid,
		(*c)["user"])

	if err != nil {
		revel.INFO.Println("ERROR: inserting into db")
		revel.INFO.Println(err)
	}

	(*c)["sid"] = newSid
}
