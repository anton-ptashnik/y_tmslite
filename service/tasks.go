package service

import (
	"errors"
	"fmt"
	"y_finalproject/persistence"
)

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

func (s *TasksService) SetTasksStatus(oldSid int64, newSid int64, pid int64) error {
	all, err := s.List(pid)
	if err != nil {
		return err
	}
	var toBeMoved []persistence.Task
	for _, v := range all {
		if v.StatusID == oldSid {
			toBeMoved = append(toBeMoved, v)
		}
	}
	var opErr []int64
	for _, v := range toBeMoved {
		v.StatusID = newSid
		err := s.Upd(v)
		if err != nil {
			opErr = append(opErr, v.ID)
		}
	}
	if len(opErr) > 0 {
		return errors.New(fmt.Sprint("failed to set status for tasks", opErr))
	}
	return nil
}
