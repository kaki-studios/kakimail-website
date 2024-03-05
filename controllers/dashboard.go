package controllers

import (
	"fmt"
	// "html/template"
	"kaki-studios/kakimail-website/auth"
	"kaki-studios/kakimail-website/user"
	"net/http"
	// "path"

	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)


func Dashboard() echo.HandlerFunc {
  return func(c echo.Context) error {
    userCookie,_ :=c.Cookie("user")

    return c.String(http.StatusOK, fmt.Sprintf("Hi, %s! you have access!", userCookie))
  }
}


// SignInForm responsible for signIn Form rendering.
func SignInForm() echo.HandlerFunc {
	return func(c echo.Context) error {
		// fp := path.Join("templates", "signIn.html")
		// tmpl, err := template.ParseFiles(fp)
		// if err != nil {
		// 	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		// }
		// if err := tmpl.Execute(c.Response().Writer, nil); err != nil {
		// 	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		// }
    return c.Render(http.StatusOK, "signIn.html", nil)

	}
}

// SignIn will be executed after SignInForm submission.
func SignIn() echo.HandlerFunc {
	return func(c echo.Context) error {
        // Load our "test" user.
		storedUser := user.LoadTestUser()
        // Initiate a new User struct.
        u := new(user.User)
        // Parse the submitted data and fill the User struct with the data from the SignIn form.
		if err := c.Bind(u); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		// Compare the stored hashed password, with the hashed version of the password that was received.
		if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(u.Password)); err != nil {
			// If the two passwords don't match, return a 401 status.
      
      return echo.NewHTTPError(http.StatusUnauthorized, "Password is incorrect")
    }
    // If password is correct, generate tokens and set cookies.
		err := auth.GenerateTokensAndSetCookies(storedUser, c)

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token is incorrect")
		}

		return c.Redirect(http.StatusMovedPermanently, "/dashboard")
	}
}
