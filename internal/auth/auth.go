package auth

import (
	"cli-todo/internal/errors"
	"cli-todo/internal/storage"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

type UserData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func Register(w http.ResponseWriter, r *http.Request) {
	var data UserData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		errors.WriteError(w, err, "fail to parse body")
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
		errors.WriteError(w, err, "fail to parse body")
		return
	}

	var storageUser storage.User
	username := storage.Db.Where("username = ?", data.Username).First(&storageUser)

	if username.Error != nil {
		msg := "failed to fetch username: " + username.Error.Error()
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

	payload := jwt.MapClaims{
		"sub": data.Username,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		errors.WriteError(w, err, "fail to sign token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	marshal, err := json.Marshal(LoginResponse{AccessToken: tokenString})
	if err != nil {
		return
	}

	w.Write(marshal)

}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil {
			http.Error(w, "Forbidden", http.StatusUnauthorized)
			fmt.Println(err)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			r.Header.Set("user_id", claims["sub"].(string))
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
