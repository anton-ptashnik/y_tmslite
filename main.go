package main

import (
	"fmt"
	"github.com/go-chi/chi"
	chiware "github.com/go-chi/chi/middleware"
	_ "github.com/joho/godotenv/autoload"
	"net/http"
	"os"
	"y_finalproject/middleware"
	"y_finalproject/persistence"
)

func main() {
	db, err := persistence.InitDb(os.Getenv("DB_CONN_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(chiware.Logger)
	r.Use(chiware.SetHeader("Content-Type", "application/json"))

	inj := NewInj()
	tasksHandler := inj.tasksHandler()
	statusesHandler := inj.statusesHandler()

	r.Route("/projects", func(r chi.Router) {
		r.Post("/", middleware.AddProject(statusesHandler.Add))
		r.Get("/", middleware.ListProjects)
		r.Get("/{id}", middleware.GetProject)
		r.Delete("/{id}", middleware.DelProject)
		r.Put("/{id}", middleware.UpdProject)
	})
	r.Route("/projects/{pid}/statuses", func(r chi.Router) {
		r.Post("/", statusesHandler.AddStatus)
		r.Get("/", statusesHandler.ListTaskStatuses)
		r.Get("/{sid}", statusesHandler.GetStatus)
		r.Delete("/{sid}", statusesHandler.DelStatus)
		r.Put("/{sid}", statusesHandler.UpdStatus)
	})
	r.Route("/projects/{pid}/tasks", func(r chi.Router) {
		r.Post("/", tasksHandler.AddTask)
		r.Get("/", tasksHandler.ListTasks)
		r.Get("/{tid}", tasksHandler.GetTask)
		r.Delete("/{tid}", tasksHandler.DelTask)
		r.Put("/{tid}", tasksHandler.UpdTask)
	})
	r.Route("/projects/{pid}/tasks/{tid}/comments", func(r chi.Router) {
		r.Post("/", middleware.AddComment)
		r.Get("/", middleware.ListComments)
		r.Get("/{cid}", middleware.GetComment)
		r.Delete("/{cid}", middleware.DelComment)
		r.Put("/{cid}", middleware.UpdComment)
	})

	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv(`PORT`)), r)
}
