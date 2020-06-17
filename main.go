package main

import (
	"github.com/go-chi/chi"
	chiware "github.com/go-chi/chi/middleware"
	"net/http"
	"y_finalproject/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(chiware.SetHeader("Content-Type", "application/json"))
	r.Route("/projects", func(r chi.Router) {
		r.Post("/", middleware.AddProject)
		r.Get("/", middleware.ListProjects)
		r.Get("/{id}", middleware.GetProject)
		r.Delete("/{id}", middleware.DelProject)
	})
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", middleware.AddTask)
		r.Get("/{id}", middleware.GetTask)
		r.Get("/", middleware.ListTasks)
		r.Delete("/{id}", middleware.DelTask)
	})
	r.Route("/tasks/{tid}/comments", func(r chi.Router) {
		r.Post("/", middleware.AddComment)
		r.Get("/", middleware.ListComments)
		r.Delete("/{cid}", middleware.DelComment)
	})
	http.ListenAndServe(":9999", r)
}
