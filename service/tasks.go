package service

import "y_finalproject/persistence"

type TasksRepo interface {
	List(pid int64) ([]persistence.Task, error)
	Del(id int64, pid int64) error
	Add(task persistence.Task) (int64, error)
	Upd(task persistence.Task) error
	Get(id int64, pid int64) (persistence.Task, error)
}

type TasksService struct {
	TasksRepo
}
