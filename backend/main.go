package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main()  {
  t := &Template{
      templates: template.Must(template.ParseGlob("../frontend/*.html")),
  }
  e := echo.New()
  e.Renderer = t;
  e.File("/", "../frontend/index.html");
  e.File("/favicon.ico", "../frontend/favicon.ico");
  e.Static("/frontend", "../frontend");
  e.GET("/:file", defaultHandler);
  e.Logger.Fatal(e.Start(":8000"))

}

//i know, its terrible!
func defaultHandler(c echo.Context) error {
  file := c.Param("file");

  log.Println("", );
  return c.Render(http.StatusOK, file+".html", nil)
}

