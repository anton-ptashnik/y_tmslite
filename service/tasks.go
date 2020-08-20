package service

import (
	"errors"
	"fmt"
	"strings"
	"y_finalproject/persistence"
)

type TasksRepo interface {
	List(pid int64) ([]persistence.Task, error)
	Del(id int64, pid int64) error
	Add(task persistence.Task) (int64, error)
	Upd(task persistence.Task) error
	Get(id int64, pid int64) (persistence.Task, error)
}

type TasksRepoTx func(tx persistence.Tx) TasksRepo

type TasksService struct {
	TasksRepo
	TasksRepoTx
	TxInitiator
}

func (s *TasksService) SetTasksStatus(oldSid int64, newSid int64, pid int64) error {
	//todo move to a separate unit
	toBeMoved, err := s.findBySid(oldSid, pid)
	if err != nil {
		return err
	}
	var failedUpdIDs []int64
	tx, err := s.TxInitiator()
	txTasks := s.TasksRepoTx(tx)
	for _, v := range toBeMoved {
		v.StatusID = newSid
		err := txTasks.Upd(v)
		if err != nil {
			failedUpdIDs = append(failedUpdIDs, v.ID)
		}
	}
	if len(failedUpdIDs) > 0 {
		tx.Rollback()
		return errors.New(fmt.Sprint("failed to set status for tasks", failedUpdIDs))
	}
	return persistence.TryCommit(tx)
}

func (s *TasksService) findBySid(sid, pid int64) ([]persistence.Task, error) {
	all, err := s.TasksRepo.List(pid)
	if err != nil {
		return nil, err
	}
	var tasks []persistence.Task
	for _, v := range all {
		if v.StatusID == sid {
			tasks = append(tasks, v)
		}
	}
	return tasks, nil
}

func (s *TasksService) List(filter TaskFilterTemplate) ([]persistence.Task, error) {
	tasks, err := s.TasksRepo.List(filter.Pid)
	if err != nil {
		return nil, err
	}

	taskFilter := newTaskFilter(filter)

	var res []persistence.Task
	for _, task := range tasks {
		if taskFilter(task) {
			res = append(res, task)
		}
	}
	return res, nil
}

type TaskFilterTemplate struct {
	Pid      int64
	Name     string
	Statuses []int64
}

func newTaskFilter(filter TaskFilterTemplate) func (task persistence.Task) bool {
	neededStatuses := map[int64]bool{}
	for _, s := range filter.Statuses {
		neededStatuses[s] = true
	}
	statusFilter := func(sid int64) bool {
		if len(neededStatuses) == 0 {
			return true
		}
		return neededStatuses[sid]
	}

	return func(task persistence.Task) bool {
		return strings.Contains(task.Name, filter.Name) && statusFilter(task.StatusID)
	}
}