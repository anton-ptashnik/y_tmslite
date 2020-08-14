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
	all, err := s.List(pid)
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