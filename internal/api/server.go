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

		r.Post("/add", addHandler)
		r.Get("/list", listHandler)
		r.Patch("/change-status", moveHandler)
		r.Put("/update", updateHandler)
		r.Delete("/delete", deleteHandler)

		r.Post("/create-board", board.CreateBoardHandler)
		r.Get("/get-board", board.GetBoardHandler)
		r.Patch("/update-board", board.UpdateBoardHandler)
		r.Delete("/delete-board", board.DeleteBoardHandler)
	})

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		return
	}
}
