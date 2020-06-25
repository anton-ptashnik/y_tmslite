package main

import (
	"github.com/go-chi/chi"
	chiware "github.com/go-chi/chi/middleware"
	_ "github.com/joho/godotenv/autoload"
	"net/http"
	"y_finalproject/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(chiware.Logger)
	r.Use(chiware.SetHeader("Content-Type", "application/json"))
	r.Route("/projects", func(r chi.Router) {
		r.Post("/", middleware.AddProject)
		r.Get("/", middleware.ListProjects)
		r.Get("/{id}", middleware.GetProject)
		r.Delete("/{id}", middleware.DelProject)
		r.Put("/{id}", middleware.UpdProject)
		r.Route("/taskstatuses", func(r chi.Router) {
			r.Post("/", middleware.AddTaskStatus)
			r.Get("/", middleware.GetTaskStatuses)
			r.Delete("/{tsid}", middleware.DelTaskStatus)
			r.Put("/{tsid}", middleware.UpdTaskStatus)
		})
	})
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", middleware.AddTask)
		r.Get("/", middleware.ListTasks)
		r.Get("/{id}", middleware.GetTask)
		r.Delete("/{id}", middleware.DelTask)
		r.Put("/{id}", middleware.UpdTask)
		r.Route("/{id}/comments", func(r chi.Router) {
			r.Post("/", middleware.AddComment)
			r.Get("/", middleware.ListComments)
			r.Delete("/{cid}", middleware.DelComment)
			r.Put("/{cid}", middleware.UpdComment)
		})
	})
	http.ListenAndServe(":9999", r)
}
