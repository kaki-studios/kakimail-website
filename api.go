package main

import (
	"database/sql"
	_ "database/sql"
	_ "encoding/base64"
	"fmt"
	_ "fmt"
	"net/http"
	_ "os"
  _ "github.com/labstack/echo-jwt/v4"

	"github.com/labstack/echo/v4"
	_ "github.com/libsql/go-libsql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
  ID int
  Name string
  Password string
}

func RegisterCreateUser(e *echo.Echo, db *sql.DB) error {
  e.POST("/create_user", func(c echo.Context) error {
    data,err := c.FormParams()
    if err != nil {
      return err
    }
    username := data.Get("username")
    password := data.Get("password")
    hash,err := bcrypt.GenerateFromPassword([]byte(password), 8);
    if err != nil {
      return err
    }
    dict := map[string]string {
      "Username": username,
      "Password": string(hash),
    }
    id := getBiggestId(db) + 1
    fmt.Println(id)
    res,err := db.Exec(fmt.Sprintf("INSERT INTO test VALUES (%d, \"%s\", \"%s\")", id, username, hash));
    if err != nil {
      fmt.Println("ERROR", err)
      return err
    }
    fmt.Println(res)

    return c.Render(http.StatusOK, "success.html", dict)
  })
  return nil
}

func getBiggestId(db *sql.DB) int32 {
  
  rows, err := db.Query("SELECT id FROM test");
  if err != nil {
    return -1
  }
  fmt.Println(rows)
  var max int32;
  for rows.Next() {
    var tmp int32;
    err = rows.Scan(&tmp)
    if tmp > max {
      max = tmp
    }
    fmt.Println(max)
    if err != nil {
      return -1
    }
  }
  return max
}
