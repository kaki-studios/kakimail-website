package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"kaki-studios/kakimail-website/auth"
	"kaki-studios/kakimail-website/controllers"
	"log"
	_ "log"
	"net/http"
	"os"

	_ "github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"

	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	_ "github.com/libsql/go-libsql"

	_ "kaki-studios/kakimail-website/user"
)

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main()  {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("hahah couldn't even load a dotenv")
  }

  dbName := fmt.Sprintf("file://%s/kakimail-website.db",os.TempDir())

  db, err := sql.Open("libsql", dbName)
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to open db %s", err)
    os.Exit(1)
  }
  
  _, err = db.Exec("CREATE TABLE IF NOT EXISTS test (id INTEGER NOT NULL, name TEXT, password TEXT, PRIMARY KEY(id), UNIQUE(id))")
	if err != nil {
		log.Fatal("no can do")
	} else {
    fmt.Println("success")
  }
  defer db.Close()


  t := &Template{
    templates: template.Must(template.ParseGlob("*/*.html")),
  }
  e := echo.New()
  e.Renderer = t;
  userGroup := e.Group("/dashboard")
  userGroup.Use(echojwt.WithConfig(echojwt.Config{
		// NewClaimsFunc:                  func(c echo.Context) jwt.Claims { return &auth.Claims{}},
    SigningKey:              []byte(auth.GetJWTSecret()),
		TokenLookup:             "cookie:access-token", // "<source>:<name>"
		ErrorHandler: auth.JWTErrorChecker,
  }))
// Attach jwt token refresher.
  userGroup.GET("", controllers.Dashboard())
  e.File("/", "static/index.html");
  e.File("/favicon.ico", "static/assets/favicon.ico");
  e.Static("/static", "static");
  e.GET("/:file", defaultHandler);
  userGroup.Use(echojwt.JWT([]byte(os.Getenv("JWT_SECRET"))))
  userGroup.Use(auth.TokenRefresherMiddleware)
  e.GET("/user/signin", controllers.SignInForm()).Name = "userSignInForm"
	e.POST("/user/signin", controllers.SignIn())
  RegisterCreateUser(e, db)

  e.Logger.Fatal(e.Start(":8000"))

}

//i know, its terrible!
func defaultHandler(c echo.Context) error {
  file := c.Param("file");

  return c.Render(http.StatusOK, fmt.Sprintf("%s.html", file), nil)
}
