package service

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"y_finalproject/persistence"
)

var errProjectNotFound = errors.New("project not found")

type statusTests struct {
	StatusesService
}

func TestStatuses(t *testing.T) {
	tests := statusTests{}
	tests.SetTasksStatusOp = func(oldSid int64, newSid int64, pid int64) error {
		return nil
	}

	t.Run("add", tests.add)
	t.Run("del last status not allowed", tests.delLastStatus)
	t.Run("status can be deleted", tests.del)
	t.Run("upd", tests.upd)
	t.Run("status seqNo can be set", tests.setSeqNo)
	t.Run("status seqNo setup fails on invalid value", tests.setInvalidSeqNo)
}

func (test *statusTests) del(t *testing.T) {
	statuses := map[int64]persistence.Status{
		1: {
			ID: 1,
			Name:  "s1",
			SeqNo: 1,
		},
		2: {
			ID: 2,
			Name:  "s2",
			SeqNo: 2,
		},
	}
	delStatusOp := func(sid, pid int64) error {
		if _, p := statuses[sid]; !p {
			return errors.New("no entry")
		}
		delete(statuses, sid)
		return nil
	}
	listStatusOp := func(pid int64) ([]persistence.Status, error) {
		var res []persistence.Status
		for _, v := range statuses {
			res = append(res, v)
		}
		return res, nil
	}
	test.StatusesRepo = statusesFakeRepo{
		delOp: delStatusOp,
		listOp: listStatusOp,
	}
	test.SetTasksStatusOp = func(oldSid int64, newSid int64, pid int64) error {
		return nil
	}
	err := test.Del(2, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, present := statuses[2]; present {
		t.Error("responded deleted but still present")
	}
}
func (test *statusTests) delLastStatus(t *testing.T) {
	status := persistence.Status{
		ID:  1,
		PID: 1,
	}
	getStatusOp := func(sid, pid int64) (persistence.Status, error) {
		return status, nil
	}
	listStatusesOp := func(pid int64) ([]persistence.Status, error) {
		return []persistence.Status{status}, nil
	}
	delStatusOp := func(id, pid int64) error {
		t.Fatal("del request sent to repo")
		return nil
	}
	test.StatusesRepo = statusesFakeRepo{
		getOp:  getStatusOp,
		listOp: listStatusesOp,
		delOp:  delStatusOp,
	}
	err := test.Del(status.ID, status.PID)
	if err == nil {
		t.Fatal("indicated last status deleted")
	}
}

func (test *statusTests) upd(t *testing.T) {
	status := persistence.Status{
		PID:  0,
		Name: "testitem",
	}
	updStatus := persistence.Status{
		Name: "updName",
	}
	test.StatusesRepo = statusesFakeRepo{
		updOp: func(_status persistence.Status) error {
			status = _status
			return nil
		},
		getOp: func(sid, pid int64) (persistence.Status, error) {
			return status, nil
		},
	}

	err := test.Upd(updStatus)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(status, updStatus) {
		t.Error(err, "expected/actual:", updStatus, status)
	}
}
func (test *statusTests) setInvalidSeqNo(t *testing.T) {
	s := persistence.Status{
		PID:   0,
		SeqNo: 1,
	}
	test.StatusesRepo = statusesFakeRepo{
		listOp: func(pid int64) ([]persistence.Status, error) {
			return []persistence.Status{
				{},
			}, nil
		},
	}

	newSeqNo := 2
	err := test.StatusesService.setSeqNo(s, newSeqNo)
	if err == nil {
		t.Fatal("invalid status seqNo accepted")
	}
}

func (test *statusTests) setSeqNo(t *testing.T) {
	// index=seqNo, ID=value
	actual := []int64{3,4,5,6}
	expectedSeq := []int64{5,3,4,6}

	test.StatusesRepo = statusesFakeRepo{
		updOp: func(status persistence.Status) error {
			actual[status.SeqNo] = status.ID
			return nil
		},
		listOp: func(pid int64) ([]persistence.Status, error) {
			var r []persistence.Status
			for i := range actual {
				r = append(r, persistence.Status{
					ID:    actual[i],
					SeqNo: i,
				})
			}
			return r, nil
		},
	}
	test.TxInitiator = func() (persistence.Tx, error) {
		return fakeTx{}, nil
	}
	test.StatusesRepoTx = func(tx persistence.Tx) StatusesRepo {
		return test.StatusesRepo
	}
	status := persistence.Status{
		ID:    5,
		SeqNo: 2,
	}
	err := test.StatusesService.setSeqNo(status, 0)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedSeq, actual) {
		t.Error("ordering mismatch. Expected/actual ([seqNo]=entryID)", expectedSeq, actual)
	}
}

func (test *statusTests) add(t *testing.T) {
	status := persistence.Status{
		ID:   2,
		Name: "testitem",
	}
	var statuses []persistence.Status

	test.StatusesRepo = statusesFakeRepo{
		addOp: func(_status persistence.Status) (int64, error) {
			statuses = append(statuses, _status)
			return _status.ID, nil
		},
		listOp: func(pid int64) ([]persistence.Status, error) {
			return statuses, nil
		},
	}
	test.TxInitiator = func() (persistence.Tx, error) {
		return fakeTx{}, nil
	}
	test.StatusesRepoTx = func(tx persistence.Tx) StatusesRepo {
		return statusesFakeRepo{}
	}
	sid, err := test.Add(status)
	if err != nil {
		t.Fatal(err)
	}
	if sid != status.ID {
		t.Error("wrong entry ID, expected/actual:", status.ID, sid)
	}

	if len(statuses) == 0 || statuses[0] != status {
		t.Error("indicated ok but no entry in a repo")
	}
}

type statusesFakeRepo struct {
	pid    int64
	items  map[int64]persistence.Status
	addOp  func(status persistence.Status) (int64, error)
	getOp  func(sid, pid int64) (persistence.Status, error)
	listOp func(pid int64) ([]persistence.Status, error)
	delOp  func(sid, pid int64) error
	updOp func(status persistence.Status) error
}

func (s statusesFakeRepo) List(pid int64) ([]persistence.Status, error) {
	return s.listOp(pid)
}

func (s statusesFakeRepo) Del(sid int64, pid int64) error {
	return s.delOp(sid, pid)
}

func (s statusesFakeRepo) Add(status persistence.Status) (int64, error) {
	return s.addOp(status)
}

func (s statusesFakeRepo) Upd(status persistence.Status) error {
	return s.updOp(status)
}

func (s statusesFakeRepo) Get(id int64, pid int64) (persistence.Status, error) {
	return s.getOp(id, pid)
}

type fakeTx struct {
}

func (f fakeTx) QueryRow(q string, args ...interface{}) *sql.Row {
	panic("implement me")
}

func (f fakeTx) Query(q string, args ...interface{}) (*sql.Rows, error) {
	panic("implement me")
}

func (f fakeTx) Exec(q string, args ...interface{}) (sql.Result, error) {
	panic("implement me")
}

func (f fakeTx) Commit() error {
	return nil
}

func (f fakeTx) Rollback() error {
	return nil
}
