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

type SetNewStatusOp func(oldSid int64, newSid int64, pid int64) error

type StatusesService struct {
	StatusesRepo
	SetNewStatusOp
}

func (s *StatusesService) Del(sid int64, pid int64) error {
	statuses, err := s.StatusesRepo.List(pid)
	if err != nil {
		return err
	}
	if len(statuses) == 1 && statuses[0].ID == sid {
		return errLastStatus
	}
	statusToBeDel, err := s.StatusesRepo.Get(sid, pid)
	if err != nil {
		return nil
	}
	var newSeqNo int
	if statusToBeDel.SeqNo == 1 {
		newSeqNo = 2
	} else {
		newSeqNo = statusToBeDel.SeqNo - 1
	}
	var targetStatusID int64
	for _, v := range statuses {
		if v.SeqNo == newSeqNo {
			targetStatusID = v.ID
			break
		}
	}
	err = s.SetNewStatusOp(sid, targetStatusID, pid)
	if err != nil {
		return err
	}
	return s.StatusesRepo.Del(sid, pid)
}

func (s *StatusesService) SetSeqNo(sid int64, pid int64, newSeqNo int) error {
	// todo refactor
	targetStatus, err := s.StatusesRepo.Get(sid, pid)
	if err != nil {
		return err
	}
	if targetStatus.SeqNo == newSeqNo {
		return nil
	}

	var di int
	var fIndex, lIndex int
	if newSeqNo > targetStatus.SeqNo {
		di = -1
		fIndex = targetStatus.SeqNo + 1
		lIndex = newSeqNo
	} else {
		di = 1
		fIndex = newSeqNo
		lIndex = targetStatus.SeqNo - 1
	}
	toBeMoved, err := s.listBySeqNo(pid, fIndex, lIndex)
	if err != nil {
		return err
	}

	for i := range toBeMoved {
		toBeMoved[i].SeqNo += di
	}
	targetStatus.SeqNo = newSeqNo
	toBeMoved = append(toBeMoved, targetStatus)

	for _, se := range toBeMoved {
		err := s.StatusesRepo.Upd(se)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *StatusesService) listBySeqNo(pid int64, l int, r int) ([]persistence.Status, error) {
	all, err := s.StatusesRepo.List(pid)
	if err != nil {
		return nil, err
	}
	if r > len(all) {
		return nil, errors.New("invalid seqNo")
	}
	var filtered []persistence.Status
	for _, s := range all {
		if s.SeqNo >= l && s.SeqNo <= r {
			filtered = append(filtered, s)
		}
	}
	return filtered, nil
}
