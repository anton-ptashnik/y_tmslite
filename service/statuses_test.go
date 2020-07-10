package service

import (
	"errors"
	"reflect"
	"testing"
	"y_finalproject/persistence"
)

var errProjectNotFound = errors.New("project not found")

type statusTests struct {
	s StatusesService
}

func TestStatuses(t *testing.T) {
	tests := statusTests{}
	tests.s.SetTasksStatusOp = func(oldSid int64, newSid int64, pid int64) error {
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
	pid := int64(0)
	r := statusesFakeRepo{
		pid: pid,
		items: map[int64]persistence.Status{},
	}
	s1 := persistence.Status{
		PID:   pid,
		Name:  "s1",
		SeqNo: 1,
	}
	s2 := persistence.Status{
		PID:   pid,
		Name:  "s2",
		SeqNo: 2,
	}
	sid, _ := r.Add(s1)
	s1.ID = sid
	sid, _ = r.Add(s2)
	s2.ID = sid
	test.s.StatusesRepo = r

	err := test.s.Del(s1.ID, pid)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.Get(s1.ID, r.pid); err == nil {
		t.Error("responded deleted but still present")
	}
}
func (test *statusTests) delLastStatus(t *testing.T) {
	r := statusesFakeRepo{
		pid:   0,
		items: map[int64]persistence.Status{},
	}
	newItemID, _ := r.Add(persistence.Status{
		PID: 0,
	})
	test.s.StatusesRepo = r
	err := test.s.Del(newItemID, 0)
	if err == nil {
		t.Fatal("last status deleted")
	}
	if _, err := r.Get(newItemID, r.pid); err != nil {
		t.Error("responded kept but is missing")
	}
}

func (test *statusTests) upd(t *testing.T) {
	status := persistence.Status{
		PID:  0,
		Name: "testitem",
	}
	r := statusesFakeRepo{
		pid:   status.ID,
		items: map[int64]persistence.Status{},
	}
	test.s.StatusesRepo = r

	newItemID, _ := r.Add(status)
	status.ID = newItemID
	status.Name += "_upd"
	err := test.s.Upd(status)
	if err != nil {
		t.Fatal(err)
	}
	actualStatus, err := r.Get(status.ID, status.PID)
	if err != nil || !reflect.DeepEqual(status, actualStatus) {
		t.Error(err, "expected/actual:", status, actualStatus)
	}
}
func (test *statusTests) setInvalidSeqNo(t *testing.T) {
	s := persistence.Status{
		PID:   0,
		SeqNo: 1,
	}
	r := statusesFakeRepo{
		pid:   s.PID,
		items: map[int64]persistence.Status{},
	}
	test.s.StatusesRepo = r

	sid, _ := r.Add(s)
	newSeqNo := 2
	err := test.s.setSeqNo(sid, s.PID, newSeqNo)
	if err == nil {
		t.Fatal("invalid status seqNo accepted")
	}
}
func (test *statusTests) setSeqNo(t *testing.T) {
	r := statusesFakeRepo{
		pid:   0,
		items: map[int64]persistence.Status{},
	}
	test.s.StatusesRepo = r

	var statuses []int64
	for n := 0; n < 5; n++ {
		id, _ := r.Add(persistence.Status{
			PID:   r.pid,
			SeqNo: n,
		})
		statuses = append(statuses, id)
	}

	newSeqNo := 3
	expected := statuses[1:4]
	iToBeMoved := expected[0]
	copy(expected[:len(expected)-1], expected[1:])
	expected[len(expected)-1] = iToBeMoved

	err := test.s.setSeqNo(iToBeMoved, r.pid, newSeqNo)
	if err != nil {
		t.Fatal(err)
	}

	actualStatuses, _ := r.List(r.pid)
	for _, s := range actualStatuses {
		if s.SeqNo >= len(statuses) || statuses[s.SeqNo] != s.ID {
			t.Error("expected/found at seqNo:",s.SeqNo, statuses[s.SeqNo], s.ID)
		}
	}

}

func (test *statusTests) add(t *testing.T) {
	status := persistence.Status{
		PID:  0,
		Name: "testitem",
		SeqNo: 1,
	}
	newStatus := persistence.Status{
		PID:  0,
		Name: "testitem_new",
		SeqNo: 1,
	}
	r := statusesFakeRepo{
		pid:   0,
		items: map[int64]persistence.Status{},
	}
	sid, _ := r.Add(status)
	status.ID = sid
	test.s.StatusesRepo = r

	sid, err := test.s.Add(newStatus)
	if err != nil {
		t.Fatal(err)
	}
	newStatus.ID = sid

	status, _ = r.Get(status.ID, status.PID)
	newStatus, _ = r.Get(newStatus.ID, newStatus.PID)
	if status.SeqNo != 2 || newStatus.SeqNo != 1 {
		t.Error("expected seqNo to be 1,2 but:", newStatus.SeqNo, status.SeqNo)
	}

}

type statusesFakeRepo struct {
	pid   int64
	items map[int64]persistence.Status
}

func (s statusesFakeRepo) UpdAll(statuses []persistence.Status) error {
	for _, v := range statuses {
		s.items[v.ID] = v
	}
	return nil
}

func (s statusesFakeRepo) List(pid int64) ([]persistence.Status, error) {
	if pid != s.pid {
		return nil, errors.New("not found")
	}
	var res []persistence.Status
	for _, s := range s.items {
		res = append(res, s)
	}
	return res, nil
}

func (s statusesFakeRepo) Del(sid int64, pid int64) error {
	if pid != s.pid {
		return errProjectNotFound
	}
	delete(s.items, sid)
	return nil
}

func (s statusesFakeRepo) Add(status persistence.Status) (int64, error) {
	if status.PID != s.pid {
		return 0, errProjectNotFound
	}
	id := int64(len(s.items) + 1)
	status.ID = id
	s.items[id] = status
	return id, nil
}

func (s statusesFakeRepo) Upd(status persistence.Status) error {
	if _, p := s.items[status.ID]; !p {
		return errors.New("not found")
	}
	s.items[status.ID] = status
	return nil
}

func (s statusesFakeRepo) Get(id int64, pid int64) (persistence.Status, error) {
	v, p := s.items[id]
	var err error
	if !p {
		err = errors.New("not found")
	}
	return v, err
}
