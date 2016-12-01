package controllers

import (
	//"github.com/jmoiron/sqlx"
	"github.com/revel/revel"
	//	"github.com/satori/go.uuid"
	"week/app"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	userName := string(c.Session["user"])
	return c.Render(userName)
}

type Comment struct {
	Id         string
	ParentPost string `db:"parent_post"`
	Author     string
	Date       string
	Content    string
	Children   []string
}

type Post struct {
	Id       string
	Title    string
	Author   string
	Date     string
	Content  string
	Comments []Comment
}

func (c App) Feed() revel.Result {
	if res := checkSession(&c.Session, &c.Flash); !res {
		return c.Redirect(User.SignIn)
	}

	Posts := []Post{}

	rows, err := app.DB.Queryx("SELECT * FROM posts ORDER BY date DESC")

	if err != nil {
		revel.INFO.Println("ERROR: querying db")
		revel.INFO.Println(err)
	} else {
		for rows.Next() {
			post := Post{}
			if err := rows.StructScan(&post); err != nil {
				revel.INFO.Println("error")
				revel.INFO.Println(err)
			} else {
				Posts = append(Posts, post)
			}
		}
	}

	defer rows.Close()

	userName := string(c.Session["user"])
	return c.Render(userName, Posts)
}

func (c App) EditPost() revel.Result {
	if res := checkSession(&c.Session, &c.Flash); !res {
		return c.Redirect(User.SignIn)
	}

	userName := string(c.Session["user"])
	return c.Render(userName)
}

func (c App) SubmitPost(titleInput, contentInput string) revel.Result {
	if res := checkSession(&c.Session, &c.Flash); !res {
		return c.Redirect(User.SignIn)
	}

	userName := string(c.Session["user"])

	res, err := app.DB.Queryx("INSERT INTO posts VALUES(DEFAULT, $1, $2, now(), $3) RETURNING id", titleInput, userName, contentInput)
	if err != nil {
		revel.INFO.Println(err)
	}
	res.Next()
	id := ""
	res.Scan(&id)
	revel.INFO.Println(id)
	return c.Redirect("/post/%s", id)
}

func (c App) SubmitComment(parentId, contentInput string) revel.Result {
	if res := checkSession(&c.Session, &c.Flash); !res {
		return c.Redirect(User.SignIn)
	}

	userName := string(c.Session["user"])

	res, err := app.DB.Queryx("INSERT INTO comments VALUES(DEFAULT, $1, $2, now(), $3) RETURNING id", parentId, userName, contentInput)
	if err != nil {
		revel.INFO.Println(err)
	}
	res.Next()
	id := ""
	res.Scan(&id)
	revel.INFO.Println(id)
	return c.Redirect("/post/%s", parentId) //TODO:
}

func (c App) ViewPost(id string) revel.Result {
	rows, err := app.DB.Queryx("SELECT * FROM posts WHERE id=$1", id)

	if err != nil {
		revel.INFO.Println("ERROR: querying db")
		c.Redirect("/")
		//		panic(err)
	}

	defer rows.Close()

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

	rows, err = app.DB.Queryx("SELECT * FROM comments WHERE parent_post=$1", id)

	if err != nil {
		revel.INFO.Println("ERROR: querying db")
		revel.INFO.Println(err)
		c.Redirect("/")
		//		panic(err)
	}

	defer rows.Close()

	comments := []Comment{}

	// 1 result guaranteed, so we don't use for
	for rows.Next() {
		comm := Comment{}
		if err := rows.StructScan(&comm); err != nil {
			revel.INFO.Println("error")
			revel.INFO.Println(err)
		}

		comments = append(comments, comm)

		if err := rows.Err(); err != nil {
			revel.INFO.Println("ERROR: in rows")
			revel.INFO.Println(err)
		}
	}

	userName := string(c.Session["user"])
	return c.Render(userName, post, comments)
}
