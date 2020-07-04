package main

import (
	"github.com/go-chi/chi"
	chiware "github.com/go-chi/chi/middleware"
	_ "github.com/joho/godotenv/autoload"
	"net/http"
	"y_finalproject/middleware"
	"y_finalproject/persistence"
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
	})
	r.Route("/projects/{pid}/statuses", func(r chi.Router) {
		r.Post("/", middleware.AddTaskStatus)
		r.Get("/", middleware.ListTaskStatuses)
		r.Get("/{sid}", middleware.GetStatus)
		r.Delete("/{sid}", middleware.DelTaskStatus)
		r.Put("/{sid}", middleware.UpdTaskStatus)
	})
	r.Route("/projects/{pid}/tasks", func(r chi.Router) {
		r.Post("/", middleware.AddTask(persistence.AddTask))
		r.Get("/", middleware.ListTasks(persistence.ListTasks))
		r.Get("/{tid}", middleware.GetTask(persistence.GetTask))
		r.Delete("/{tid}", middleware.DelTask(persistence.DelTask))
		r.Put("/{tid}", middleware.UpdTask(persistence.UpdTask))
	})
	r.Route("/projects/{pid}/tasks/{tid}/comments", func(r chi.Router) {
		r.Post("/", middleware.AddComment)
		r.Get("/", middleware.ListComments)
		r.Delete("/{cid}", middleware.DelComment)
		r.Put("/{cid}", middleware.UpdComment)
	})

	http.ListenAndServe(":9999", r)
}
