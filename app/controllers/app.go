package controllers

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/revel/revel"
	"strings"
	"week/app"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	if res := checkSession(&c.Session, &c.Flash); res {
		return c.Redirect("/feed")
	}
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
	Id      string
	Title   string
	Author  string
	Date    string
	Content string
	Tags    []string
}

func (post *Post) scan(rows *sqlx.Rows) error {
	return rows.Scan(&post.Id,
		&post.Title,
		&post.Author,
		&post.Date,
		&post.Content,
		pq.Array(&post.Tags))
}

func (c App) Feed() revel.Result {
	if res := checkSession(&c.Session, &c.Flash); !res {
		return c.Redirect(User.SignIn)
	}

	posts := []Post{}
	placeholder := []string{"mate", "not technology"}

	// pull this user's tags from db
	rows, err := app.DB.Queryx("SELECT * FROM posts WHERE tags && $1 ORDER BY date DESC", pq.Array(&placeholder))

	if err != nil {
		revel.INFO.Println("ERROR: querying db")
		revel.INFO.Println(err)
	} else {
		for rows.Next() {
			post := Post{}
			if err := post.scan(rows); err != nil {
				revel.INFO.Println("error")
				revel.INFO.Println(err)
			} else {
				posts = append(posts, post)
			}
		}
	}

	defer rows.Close()

	userName := string(c.Session["user"])
	return c.Render(userName, posts)
}

func (c App) EditPost() revel.Result {
	if res := checkSession(&c.Session, &c.Flash); !res {
		return c.Redirect(User.SignIn)
	}

	userName := string(c.Session["user"])
	return c.Render(userName)
}

func (c App) SubmitPost(titleInput, tagsInput, contentInput string) revel.Result {
	if res := checkSession(&c.Session, &c.Flash); !res {
		return c.Redirect(User.SignIn)
	}
	tags := strings.Split(tagsInput, ",")

	if len(tags) > 3 {

		c.Flash.Error("You can have a maximum of 3 tags")
		return c.Redirect(App.EditPost)
	}

	trimmedTags := []string{}
	for _, tag := range tags {
		trimmedTags = append(trimmedTags, strings.Trim(tag, " "))
	}

	revel.INFO.Println(tags)

	userName := string(c.Session["user"])

	res, err := app.DB.Queryx("INSERT INTO posts VALUES(DEFAULT, $1, $2, now(), $3, $4) RETURNING id",
		titleInput,
		userName,
		contentInput,
		pq.Array(trimmedTags))

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
	return c.Redirect("/post/%s?commentId=%s", parentId, id)
}

func (c App) ViewPost(id string, commentId string) revel.Result {
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

	if err := post.scan(rows); err != nil {

		revel.INFO.Println("error")
		revel.INFO.Println(err)
	}

	revel.INFO.Println("tags: ", post.Tags)

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
	return c.Render(userName, post, comments, commentId)
}

func (c App) ViewComment(id string) revel.Result {
	rows, err := app.DB.Queryx("SELECT * FROM comments WHERE id=$1", id)

	if err != nil {
		revel.INFO.Println("ERROR: querying db")
		c.Redirect("/")
	}

	defer rows.Close()

	comment := Comment{}

	// 1 result guaranteed, so we don't use for
	rows.Next()

	if err := rows.StructScan(&comment); err != nil {
		revel.INFO.Println("error")
		revel.INFO.Println(err)
	}

	if err := rows.Err(); err != nil {
		revel.INFO.Println("ERROR: in rows")
		revel.INFO.Println(err)
	}

	userName := string(c.Session["user"])
	return c.Render(comment, userName)
}

func (c App) BrowseTag(tag string) revel.Result {
	posts := []Post{}
	tags := []string{tag}

	rows, err := app.DB.Queryx("SELECT * FROM posts WHERE tags && $1 ORDER BY date DESC", pq.Array(&tags))

	if err != nil {
		revel.INFO.Println("ERROR: querying db")
		revel.INFO.Println(err)
	} else {
		for rows.Next() {
			post := Post{}
			if err := post.scan(rows); err != nil {
				revel.INFO.Println("error")
				revel.INFO.Println(err)
			} else {
				posts = append(posts, post)
			}
		}
	}

	defer rows.Close()

	userName := string(c.Session["user"])
	return c.Render(userName, posts)
}
