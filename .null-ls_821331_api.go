package main

import (
	_ "fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/libsql/go-libsql"
)


func CreateUser(c echo.Context) error {
  data,err := c.FormParams();
  if err != nil {
    return err
  }
  username := data.Get("username")
  dict := map[string]string {
    "Username": username,
  }

  return c.Render(http.StatusOK, "success.html", dict)
}
