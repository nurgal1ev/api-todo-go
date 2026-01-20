package main

import (
	"cli-todo/internal/api"
	"cli-todo/internal/storage"
)

func main() {
	storage.NewDB()
	api.HTTPServer()
}
