package main

import (
	"database/sql"
	_ "database/sql"
	"encoding/base64"
	"fmt"
	_ "fmt"
	"math/rand"
	"net/http"
	_ "os"

	"github.com/labstack/echo/v4"
	_ "github.com/libsql/go-libsql"
)

type User struct {
  ID int
  Name string
  Passhash string
}

func RegisterCreateUser(e *echo.Echo, db *sql.DB) error {
  
  



  e.POST("/create_user", func(c echo.Context) error {
    data,err := c.FormParams()
    if err != nil {
      return err
    }
    username := data.Get("username")
    password := data.Get("password")
    dst := base64.RawStdEncoding.EncodeToString([]byte(password))
    dict := map[string]string {
      "Username": username,
      "Password": dst,
    }
    res,err := db.Exec(fmt.Sprintf("INSERT INTO test VALUES (%d, \"%s\", \"%s\")", rand.Int31() / 1000, username, dst));
    if err != nil {
      fmt.Println(err)
      return err
    }
    fmt.Println(res)

    return c.Render(http.StatusOK, "success.html", dict)
  })
  return nil
}
