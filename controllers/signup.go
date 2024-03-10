package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"kaki-studios/kakimail-website/auth"
	"kaki-studios/kakimail-website/user"
)

// e.GET("/user/signin", controllers.SignInForm()).Name = "userSignInForm"
// e.POST("/user/signin", controllers.SignIn())

func SignUpForm() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "signup.html", nil)
	}
}

func SignUp(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(user.User)
		// Parse the submitted data and fill the User struct
		if err := c.Bind(u); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
		if err != nil {
			return err
		}
		// this is safe, using ?
		res, err := db.Exec(
			"INSERT INTO users (id, name, password) VALUES (NULL, ?, ?)", u.Name, hash,
		)
		if err != nil {
			fmt.Println("ERROR", err)
			return err
		}
		val, err := res.RowsAffected()
		if err != nil {
			fmt.Println("ERROR", err)
			return err
		} else if val == 0 {
			return c.String(400, "user already exists")
		}

		err = auth.GenerateTokensAndSetCookies(u, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token is incorrect")
		}

		return c.Redirect(http.StatusMovedPermanently, "/dashboard")
		// return c.Render(http.StatusOK, "success.html", dict)
	}
}
