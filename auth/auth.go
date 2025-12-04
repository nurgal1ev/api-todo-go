package auth

import (
	"cli-todo/storage"
	"log"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(login, hashedPassword) *User {

}
