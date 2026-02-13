package api

import (
	"api-todo-go/internal/auth"
	"api-todo-go/internal/board"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func HTTPServer() {
	r := chi.NewRouter()

	r.Get("/auth/login", auth.Login)
	r.Post("/auth/register", auth.Register)

	r.Group(func(r chi.Router) {

		r.Use(auth.AuthMiddleware)
		r.Post("/add", addHandler)
		r.Get("/list", listHandler)
		r.Patch("/change-status", moveHandler)
		r.Put("/update", updateHandler)
		r.Delete("/delete", deleteHandler)

		r.Route("/board", func(r chi.Router) {
			r.Post("/create", board.CreateBoardHandler)
			r.Get("/get", board.GetBoardHandler)
			r.Patch("/update", board.UpdateBoardHandler)
			r.Delete("/delete", board.DeleteBoardHandler)
			r.Post("/invite", board.InviteUserHandler)
		})
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
