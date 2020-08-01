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
	r := newStatusesRepo(nil)
	s := service.StatusesService{r, i.ts.SetTasksStatus, newStatusesRepo, newTx}
	return middleware.StatusesHandler{s}
}

func (i *inj) tasksHandler() middleware.TasksHandler {
	return middleware.TasksHandler{i.ts}
}

func newStatusesRepo(tx persistence.Tx) service.StatusesRepo {
	r := persistence.NewStatusesRepo(tx)
	return r
}

func newTx() (persistence.Tx,error) {
	tx, err := persistence.NewTx()
	return tx, err
}