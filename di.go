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
	s := service.TasksService{&r, newTasksRepo, newTx}
	return inj{s}
}

func (i *inj) projectsHandler() middleware.ProjectsHandler {
	s := newProjectsService()
	return middleware.ProjectsHandler{s}
}

func (i *inj) statusesHandler() middleware.StatusesHandler {
	r := newStatusesRepo(nil)
	s := service.StatusesService{r, i.ts.SetTasksStatus, newStatusesRepo, newTx}
	return middleware.StatusesHandler{&s}
}

func (i *inj) tasksHandler() middleware.TasksHandler {
	return middleware.TasksHandler{i.ts}
}

func newStatusesRepo(tx persistence.Tx) service.StatusesRepo {
	return persistence.NewStatusesRepo(tx)
}

func newTasksRepo(tx persistence.Tx) service.TasksRepo {
	return persistence.NewTasksRepo(tx)
}

func newProjectsRepo(tx persistence.Tx) service.ProjectsRepo {
	return persistence.NewProjectsRepo(tx)
}

func newTx() (persistence.Tx,error) {
	tx, err := persistence.NewTx()
	return tx, err
}

func newProjectsService() middleware.ProjectsService {
	return &service.ProjectsService{
		ProjectsRepo:  newProjectsRepo(nil) ,
		ProjectsRepoTx: newProjectsRepo,
		TxInitiator:    newTx,
		InsertStatusOpTx: service.InsertStatusOpTx(func(tx persistence.Tx) service.InsertStatusOp {
			return newStatusesRepo(tx).Add
		}),
	}
}