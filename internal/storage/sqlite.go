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

	Tasks []Task
}

type Board struct {
	gorm.Model
	Name        string
	Description string
	Status      string

	UserID uint
	User   User
}
type Task struct {
	gorm.Model
	Text string
	Done bool

	UserID uint
	User   User
}

func NewDB() {
	var err error
	Db, err = gorm.Open(sqlite.Open("./storage.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	Db.AutoMigrate(&User{}, &Task{}, &Board{})
}
