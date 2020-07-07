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
	r := persistence.StatusesRepo{}
	s := service.StatusesService{&r, i.ts.SetTasksStatus}
	return middleware.StatusesHandler{s}
}

func (i *inj) tasksHandler() middleware.TasksHandler {
	return middleware.TasksHandler{i.ts}
}
