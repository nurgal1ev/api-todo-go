package main

import (
	"api-todo-go/internal/api"
	"api-todo-go/internal/storage"
)

func main() {
	storage.NewDB()
	api.HTTPServer()
}
