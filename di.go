package main

import (
	"y_finalproject/middleware"
	"y_finalproject/persistence"
	"y_finalproject/service"
)

type inj struct {
	ts service.TasksService
}

func NewInj() inj {
	r := persistence.TasksRepo{}
	s := service.TasksService{&r}
	return inj{s}
}
func (i *inj) statusesHandler() middleware.StatusesHandler {
	r := NewStatusesRepo(nil)
	s := service.StatusesService{r, i.ts.SetTasksStatus,NewStatusesRepo}
	return middleware.StatusesHandler{s}
}

func (i *inj) tasksHandler() middleware.TasksHandler {
	return middleware.TasksHandler{i.ts}
}

func NewStatusesRepo(tx *persistence.Tx) service.StatusesRepo {
	r := persistence.NewStatusesRepo(tx)
	return r
}
