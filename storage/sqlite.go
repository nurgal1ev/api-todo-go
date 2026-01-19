package storage

import (
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var Db *gorm.DB

type User struct {
	gorm.Model
	Username     string
	Email        string
	PasswordHash string
	Board        []Board
}

type Board struct {
	gorm.Model
	Name   string
	UserID uint
	User   User
	Tasks  []Task
}
type Task struct {
	gorm.Model
	Text    string
	BoardID uint
	Board   Board
}

func NewDB() {
	var err error
	Db, err = gorm.Open(sqlite.Open("./storage.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	Db.AutoMigrate(&User{}, &Board{}, &Task{})
}
