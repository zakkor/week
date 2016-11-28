package controllers

import (
	//	"encoding/hex"
	"github.com/revel/revel"
	//	"golang.org/x/crypto/scrypt"
	//	"math/rand"
	//	"strconv"
	//	"time"
	//	"week/app"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	userName := string(c.Session["user"])
	return c.Render(userName)
}

func (c App) Hello(myName string) revel.Result {
	return c.Render(myName)
}
