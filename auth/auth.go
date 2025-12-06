package auth

import (
	"cli-todo/storage"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var data UserData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		msg := "fail to write HTTP response: " + err.Error()
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response: " + err.Error())
			return
		}
		fmt.Println(msg)
		return
	}

	user := storage.User{
		Username:     data.Username,
		Email:        data.Email,
		PasswordHash: data.Password,
	}

	result := storage.Db.Create(&user)
	if result.Error != nil {
		msg := "fail to write HTTP response: " + result.Error.Error()
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response: " + err.Error())
			return
		}
		fmt.Println(msg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("successful registration"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	var data UserData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		msg := "fail to write HTTP response: " + err.Error()
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response: " + err.Error())
			return
		}
		fmt.Println(msg)
		return
	}

	var storageUser storage.User
	username := storage.Db.Where("username = ?", data.Username).First(&storageUser)

	if username.Error != nil {
		msg := "fail to write HTTP response: " + username.Error.Error()
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response: " + err.Error())
			return
		}
		fmt.Println(msg)
		return
	}

	if storageUser.PasswordHash != data.Password {
		msg := "wrong password"
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response: " + err.Error())
			return
		}
		fmt.Println(msg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successful login"))
}
