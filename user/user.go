package user

import "golang.org/x/crypto/bcrypt"


type User struct {
  Password string `json:"password" form:"password"`
  Name string `json:"name" form:"name"`
}

func LoadTestUser() *User {
  //demo
  hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test"), 8);
  return &User{Password: string(hashedPassword), Name: "Test user"}
}
