package api

import (
	"api-todo-go/internal/auth"
	"api-todo-go/internal/board"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func HTTPServer() {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		r.Post("/tasks", addHandler)
		r.Get("/tasks", listHandler)
		r.Patch("/tasks/{id}/status", moveHandler)
		r.Put("/tasks/{id}", updateHandler)
		r.Delete("/tasks/{id}", deleteHandler)

		r.Post("/boards", board.CreateBoardHandler)
		r.Get("/boards/{id}", board.GetBoardHandler)
		r.Patch("/boards/{id}", board.UpdateBoardHandler)
		r.Delete("/boards/{id}", board.DeleteBoardHandler)
	})

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		return
	}
}
