package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	_ "log"
	"net/http"
	"os"

	_ "github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	_ "github.com/libsql/go-libsql"

	"kaki-studios/kakimail-website/auth"
	"kaki-studios/kakimail-website/controllers"
	_ "kaki-studios/kakimail-website/user"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// TODO: https: https://echo.labstack.com/docs/cookbook/auto-tls#server
	err := godotenv.Load()
	if err != nil {
		log.Fatal("hahah couldn't even load a .env file")
	}
	// OPEN A DB CONNECTION
	dbName := fmt.Sprintf("file://%s/kakimail.db", os.TempDir())
	db, err := sql.Open("libsql", dbName)
	if err != nil {
		fmt.Printf("failed to open db: %s", err)
		os.Exit(1)
	} else {
		fmt.Println("successful database connection")
	}
	// create the table if not exists
	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT UNIQUE, password TEXT)",
	)
	if err != nil {
		log.Fatal("no can do", err)
	}
	// we need some indices yk
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS users_name ON users(name);")
	if err != nil {
		log.Fatal("no can do", err)
	}
	// never forget
	defer db.Close()
	// now, set up the web server
	t := &Template{
		templates: template.Must(template.ParseGlob("*/*.html")),
	}
	e := echo.New()
	e.Renderer = t
	userGroup := e.Group("/dashboard")
	userGroup.Use(echojwt.WithConfig(echojwt.Config{
		// NewClaimsFunc:                  func(c echo.Context) jwt.Claims { return &auth.Claims{}},
		SigningKey:   []byte(auth.GetJWTSecret()),
		TokenLookup:  "cookie:access-token", // "<source>:<name>"
		ErrorHandler: auth.JWTErrorChecker,
	}))
	userGroup.GET("", controllers.Dashboard())
	e.File("/", "static/index.html")
	e.File("/favicon.ico", "static/assets/favicon.ico")
	e.Static("/static", "static")
	e.GET("/:file", defaultHandler)
	userGroup.Use(echojwt.JWT([]byte(os.Getenv("JWT_SECRET"))))
	// Attach jwt token refresher.
	userGroup.Use(auth.TokenRefresherMiddleware)
	e.GET("/user/signin", controllers.SignInForm()).Name = "userSignInForm"
	e.POST("/user/signin", controllers.SignIn(db))

	e.GET("/user/signup", controllers.SignUpForm()).Name = "userSignUpForm"
	e.POST("/user/signup", controllers.SignUp(db))
	e.Logger.Fatal(e.Start(":8000"))
}

// i know, its terrible!
func defaultHandler(c echo.Context) error {
	file := c.Param("file")

	return c.Render(http.StatusOK, fmt.Sprintf("%s.html", file), nil)
}
