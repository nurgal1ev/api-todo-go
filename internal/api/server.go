package api

import (
	"api-todo-go/internal/auth"
	"api-todo-go/internal/board"
	"fmt"
	"net/http"
)

func HTTPServer() {
	router := http.NewServeMux()
	router.Handle("/add", auth.AuthMiddleware(http.HandlerFunc(addHandler)))
	router.Handle("/list", auth.AuthMiddleware(http.HandlerFunc(listHandler)))
	router.Handle("/done", auth.AuthMiddleware(http.HandlerFunc(doneHandler)))
	router.Handle("/delete", auth.AuthMiddleware(http.HandlerFunc(deleteHandler)))
	router.Handle("/update", auth.AuthMiddleware(http.HandlerFunc(updateHandler)))
	router.Handle("/create-board", auth.AuthMiddleware(http.HandlerFunc(board.CreateBoardHandler)))
	router.Handle("/delete-board", auth.AuthMiddleware(http.HandlerFunc(board.DeleteBoardHandler)))
	router.Handle("/update-board", auth.AuthMiddleware(http.HandlerFunc(board.UpdateBoardHandler)))
	router.Handle("/get-board", auth.AuthMiddleware(http.HandlerFunc(board.GetBoardHandler)))
	router.HandleFunc("/auth/register", auth.Register)
	router.HandleFunc("/auth/login", auth.Login)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("http server started")

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("HTTP server error", err)
		return
	}
}
