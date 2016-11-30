package controllers

import (
	"encoding/hex"
	"github.com/revel/revel"
	"golang.org/x/crypto/scrypt"

	"week/app"
)

type User struct {
	*revel.Controller
}

type UserModel struct {
	Email, Name, Password, Session string
}

const salt string = "($#*shorseyAJSKLDSJAKLS914182901skj)"

func (c *User) checkErrors() bool {
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return true
	}
	return false
}

func hash(toHash string) string {
	result, err := scrypt.Key([]byte(toHash), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		revel.INFO.Fatal("Can't hash")
		panic("error hashing")
	}
	return hex.EncodeToString(result)
}

func (c User) SignIn() revel.Result {
	//	revel.INFO.Println()
	userName := string(c.Session["user"])
	return c.Render(userName)
}

func (c *User) SubmitSignIn(signUsername, signPassword string, rememberMe bool) revel.Result {
	c.Validation.Required(signUsername).Message("Your user name is required!")
	c.Validation.Required(signPassword).Message("Your password is required!")
	c.Validation.MinSize(signPassword, 6).Message("Your password is not long enough!")

	// first check if the user has input the needed fields
	// just return early if he didnt
	if c.checkErrors() {
		return c.Redirect(User.SignIn)
	}

	foundUser := false
	passwordMatches := false

	rows, err := app.DB.Queryx("SELECT * FROM users WHERE name=$1", signUsername)

	userModel := UserModel{}
	if err != nil {
		revel.INFO.Println("ERROR: querying db")
		revel.INFO.Println(err)

	} else {
		rows.Next()

		if err := rows.StructScan(&userModel); err != nil {
			revel.INFO.Println("error")
			revel.INFO.Println(err)
		}
		foundUser = true

		if userModel.Password == hash(signPassword) {
			passwordMatches = true
		}
		if err := rows.Err(); err != nil {
			revel.INFO.Println("ERROR: in rows")
			revel.INFO.Println(err)
		}
	}

	defer rows.Close()

	if !foundUser || !passwordMatches {
		c.Validation.Error("Wrong username or password!")
	}

	// check if username was found
	// and if password matches
	if c.checkErrors() {
		return c.Redirect(User.SignIn)
	}

	revel.INFO.Println(revel.Slug("The impact of globalization"))

	// all good, we log the user in and redirect him to Index
	c.Session["user"] = signUsername
	generateSession(&c.Session)

	c.Flash.Success("Welcome", signUsername)
	return c.Redirect(App.Feed)
}

func (c *User) Register(registerEmail, registerUsername,
	registerPassword, registerConfirm string) revel.Result {

	c.Validation.Required(registerEmail).Message("Your email is required!")
	c.Validation.Required(registerUsername).Message("Your user name is required!")
	c.Validation.Required(registerPassword).Message("Your password is required!")
	c.Validation.Required(registerConfirm).Message("Your password confirmation is required!")
	c.Validation.MinSize(registerPassword, 6).Message("Your password is not long enough!")

	if c.checkErrors() {
		return c.Redirect(User.SignIn)
	}

	if registerPassword != registerConfirm {
		c.Validation.Error("Passwords do not match")
	}

	if c.checkErrors() {
		return c.Redirect(User.SignIn)
	}

	_, err := app.DB.Exec("INSERT INTO users VALUES ($1, $2, $3);",
		registerEmail,
		registerUsername,
		hash(registerPassword))

	if err != nil {
		revel.INFO.Println("ERROR: inserting into db")
		revel.INFO.Println(err)
	}

	// all good, we log the user in and redirect him to Index
	c.Session["user"] = registerUsername
	generateSession(&c.Session)
	c.Flash.Success("Welcome", registerUsername)

	return c.Redirect(App.Feed)
}

func (c *User) Logout() revel.Result {
	if res := checkSession(&c.Session, &c.Flash); !res {
		return c.Redirect(User.SignIn)
	}

	c.Session = make(map[string]string)
	return c.Redirect(App.Index)
}

func (c *User) Settings() revel.Result {
	if res := checkSession(&c.Session, &c.Flash); !res {
		return c.Redirect(User.SignIn)
	}

	userName := string(c.Session["user"])
	return c.Render(userName)
}
