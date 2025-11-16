package main

import (
	"cli-todo/api"
	"cli-todo/storage"
)

func main() {
	storage.NewDB()
	api.HTTPServer()
}
