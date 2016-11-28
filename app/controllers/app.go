package controllers

import (
	//"github.com/jmoiron/sqlx"
	"github.com/revel/revel"
	"week/app"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	userName := string(c.Session["user"])
	return c.Render(userName)
}

type Post struct {
	Id      string
	Title   string
	Author  string
	Date    string
	Picture string
	Content string
}

func (c App) Feed() revel.Result {
	if res := checkSession(&c.Session, &c.Flash); !res {
		return c.Redirect(User.SignIn)
	}

	Posts := []Post{Post{
		"asdasdas",
		"The impact of globalization",
		"John",
		"12",
		"asdf",
		"asdasdasdasdasd",
	}}

	userName := string(c.Session["user"])
	return c.Render(userName, Posts)
}

func (c App) ViewPost(id string) revel.Result {
	rows, err := app.DB.Queryx("SELECT * FROM posts WHERE id=$1", id)
	defer rows.Close()

	if err != nil {
		revel.INFO.Println("ERROR: querying db")
		panic(err)
	}

	post := Post{}

	// 1 result guaranteed, so we don't use for
	rows.Next()

	if err := rows.StructScan(&post); err != nil {
		revel.INFO.Println("error")
		revel.INFO.Println(err)
	}

	if err := rows.Err(); err != nil {
		revel.INFO.Println("ERROR: in rows")
		revel.INFO.Println(err)
	}

	userName := string(c.Session["user"])
	return c.Render(userName, post)
}
