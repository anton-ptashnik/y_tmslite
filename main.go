package main

import (
	"github.com/go-chi/chi"
	"net/http"
	"y_finalproject/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Route("/projects", func(r chi.Router) {
		r.Post("/", middleware.AddProject)
		r.Get("/", middleware.ListProjects)
		r.Delete("/{id}", middleware.DelProject)
	})
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", middleware.AddTask)
		r.Get("/", middleware.ListTasks)
		r.Delete("/{id}", middleware.DelTask)
	})
	http.ListenAndServe(":9999", r)
}
