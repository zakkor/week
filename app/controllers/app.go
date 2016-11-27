package controllers

import (
	"encoding/hex"

	"github.com/revel/revel"
	"golang.org/x/crypto/scrypt"
	"week/app"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Hello(myName string) revel.Result {
	return c.Render(myName)
}

func (c App) SignIn() revel.Result {
	return c.Render()
}

func (c *App) checkErrors() revel.Result {
	if c.Validation.HasErrors() {

		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.SignIn)
	}
	return nil
}

const salt string = "($#*shorseyAJSKLDSJAKLS914182901skj)"

func (c *App) SubmitSignIn(signUsername, signPassword string, rememberMe bool) revel.Result {
	c.Validation.Required(signUsername).Message("Your user name is required!")
	c.Validation.Required(signPassword).Message("Your password is required!")
	c.Validation.MinSize(signPassword, 6).Message("Your password is not long enough!")

	var result revel.Result

	// first check if the user has input the needed fields
	// just return early if he didnt
	result = c.checkErrors()
	if result != nil {
		return result
	}

	rows, err := app.DB.Query("SELECT * FROM users WHERE name=$1", signUsername)
	defer rows.Close()

	foundUser := false
	passwordMatches := false

	if err != nil {
		revel.INFO.Println("ERROR: querying db")
		revel.INFO.Println(err)
	} else {
		for rows.Next() {
			var email, name, password string
			if err := rows.Scan(&email, &name, &password); err != nil {
				revel.INFO.Println("error")
			}
			foundUser = true

			//hash password and compare
			hashedPass, err := scrypt.Key([]byte(signPassword), []byte(salt), 16384, 8, 1, 32)
			if err != nil {
				revel.INFO.Println("ERROR: hashing")
				revel.INFO.Println(err)
			}

			if password == hex.EncodeToString(hashedPass) {
				passwordMatches = true
			}
		}
		if err := rows.Err(); err != nil {
			revel.INFO.Println("ERROR: in rows")
			revel.INFO.Println(err)
		}
	}

	if !foundUser || !passwordMatches {
		c.Validation.Error("Wrong username or password!")
	}

	// check if username was found
	// and if password matches
	result = c.checkErrors()
	if result != nil {
		return result
	}

	return c.Redirect(App.Index)
}

func (c *App) Register(registerEmail, registerUsername,
	registerPassword, registerConfirm string) revel.Result {

	c.Validation.Required(registerEmail).Message("Your email is required!")
	c.Validation.Required(registerUsername).Message("Your user name is required!")
	c.Validation.Required(registerPassword).Message("Your password is required!")
	c.Validation.Required(registerConfirm).Message("Your password confirmation is required!")
	c.Validation.MinSize(registerPassword, 6).Message("Your password is not long enough!")

	var result revel.Result
	result = c.checkErrors()
	if result != nil {
		return result
	}

	if registerPassword != registerConfirm {
		c.Validation.Error("Passwords do not match")
	}
	result = c.checkErrors()
	if result != nil {
		return result
	}

	// hash and add to db
	hashedPass, err := scrypt.Key([]byte(registerPassword), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		revel.INFO.Println("ERROR: hashing")
		revel.INFO.Println(err)
	}

	_, err = app.DB.Exec("INSERT INTO users VALUES ($1, $2, $3);",
		registerEmail,
		registerUsername,
		hex.EncodeToString(hashedPass))

	if err != nil {
		revel.INFO.Println("ERROR: inserting into db")
		revel.INFO.Println(err)
	}

	return c.Redirect(App.Index)
}
