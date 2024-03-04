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

	"github.com/labstack/echo/v4"
	_ "github.com/libsql/go-libsql"
)

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main()  {

  dbName := fmt.Sprintf("file://%s/kakimail-website.db",os.TempDir())

  db, err := sql.Open("libsql", dbName)
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to open db %s", err)
    os.Exit(1)
  }
  _, err = db.Exec("CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY, name TEXT, passhash TEXT)")
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
  e.File("/", "static/index.html");
  e.File("/favicon.ico", "static/assets/favicon.ico");
  e.Static("/static", "static");
  e.GET("/:file", defaultHandler);
  RegisterCreateUser(e, db)

  e.Logger.Fatal(e.Start(":8000"))

}

//i know, its terrible!
func defaultHandler(c echo.Context) error {
  file := c.Param("file");

  return c.Render(http.StatusOK, fmt.Sprintf("%s.html", file), nil)
}
