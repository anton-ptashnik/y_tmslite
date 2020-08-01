package service

import (
	"errors"
	"fmt"
	"y_finalproject/persistence"
)

type StatusesRepo interface {
	List(pid int64) ([]persistence.Status, error)
	Del(id int64, pid int64) error
	Add(persistence.Status) (int64, error)
	Upd(persistence.Status) error
	Get(id int64, pid int64) (persistence.Status, error)
}

type SetTasksStatusOp func(oldSid int64, newSid int64, pid int64) error
type StatusesRepoTx func(tx persistence.Tx) StatusesRepo
type TxInitiator func() (persistence.Tx, error)

type StatusesService struct {
	StatusesRepo
	SetTasksStatusOp
	StatusesRepoTx
	TxInitiator
}

func (s *StatusesService) Del(sid int64, pid int64) error {
	statuses, err := s.StatusesRepo.List(pid)
	if err != nil {
		return err
	}
	if len(statuses) == 1 && statuses[0].ID == sid {
		return errLastStatusDelAttempt
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
	err = s.SetTasksStatusOp(sid, targetStatusID, pid)
	if err != nil {
		return err
	}
	err = s.setSeqNo(sid, pid, len(statuses))
	if err != nil {
		return err
	}
	return s.StatusesRepo.Del(sid, pid)
}

func (s *StatusesService) setSeqNo(sid int64, pid int64, newSeqNo int) error {
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
	err = s.moveStatuses(toBeMoved, di)
	if err != nil {
		return err
	}
	targetStatus.SeqNo = newSeqNo
	return s.StatusesRepo.Upd(targetStatus)
}

func (s *StatusesService) filterBySeqNo(all []persistence.Status, l int, r int) ([]persistence.Status, error) {
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

func (s *StatusesService) Upd(upd persistence.Status) error {
	current, err := s.Get(upd.ID, upd.PID)
	if err != nil {
		return err
	}
	if current.SeqNo != upd.SeqNo {
		err := s.setSeqNo(upd.ID, upd.PID, upd.SeqNo)
		if err != nil {
			return err
		}
	}
	return s.StatusesRepo.Upd(upd)
}

func (s *StatusesService) Add(status persistence.Status) (int64, error) {
	all, err := s.StatusesRepo.List(status.PID)
	toBeMoved, err := s.filterBySeqNo(all, status.SeqNo, len(all))
	if err != nil {
		return 0, err
	}
	err = s.moveStatuses(toBeMoved, 1)
	if err != nil {
		return 0,err
	}
	return s.StatusesRepo.Add(status)
}

func (s *StatusesService) moveStatuses(toBeMoved []persistence.Status, d int) error {
	for i := range toBeMoved {
		toBeMoved[i].SeqNo += d
	}

	tx, err := s.TxInitiator()
	if err != nil {
		return fmt.Errorf("err in Tx init: %s", err)
	}
	txRepo := s.StatusesRepoTx(tx)
	for _, status := range toBeMoved {
		err := txRepo.Upd(status)
		if err != nil {
			errR := tx.Rollback()
			return errors.New(fmt.Sprint("failed to change seqNo. cause err/rollback err:", err, errR))
		}
	}
	return persistence.TryCommit(tx)
}