package service

import (
	"errors"
	"y_finalproject/persistence"
)

var (
	errLastStatus = errors.New("last status cannot be deleted")
)

type StatusesRepo interface {
	List(pid int64) ([]persistence.Status, error)
	Del(id int64, pid int64) error
	Add(persistence.Status) (int64, error)
	Upd(persistence.Status) error
	Get(id int64, pid int64) (persistence.Status, error)
}

type StatusesService struct {
	StatusesRepo
}

func (s *StatusesService) Del(sid int64, pid int64) error {
	statuses, err := s.StatusesRepo.List(pid)
	if err != nil {
		return err
	}
	if len(statuses) == 1 && statuses[0].ID == sid {
		return errLastStatus
	}
	return s.StatusesRepo.Del(sid, pid)
}
