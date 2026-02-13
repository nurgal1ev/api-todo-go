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

type Status struct {
	gorm.Model
	Name  string
	Tasks []Task

	BoardID uint
	Board   Board
}
type Board struct {
	gorm.Model
	Name        string
	Description string
	CreatorID   uint
	Creator     User
}
type Task struct {
	gorm.Model
	Text     string
	StatusID uint
	Status   Status
	UserID   uint
	User     User
}

type UserBoards struct {
	gorm.Model
	BoardID uint
	UserID  uint
	Role    string
}

func NewDB() {
	var err error
	Db, err = gorm.Open(sqlite.Open("./storage.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = Db.AutoMigrate(&User{}, &Task{}, &Board{}, &Status{})
	if err != nil {
		return
	}
}
