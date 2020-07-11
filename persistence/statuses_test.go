package persistence

import (
	"fmt"
	"reflect"
	"testing"
)

type statusesTests struct {
	StatusesRepo
}

func TestStatuses(t *testing.T) {
	projects, err := prepareProjects(db, 1)
	if err != nil {
		t.Fatal(fmt.Sprint("projects prep failed:", err))
	}
	project := projects[0]
	statuses, err := prepareStatuses(db, project, 2)
	if err != nil {
		t.Fatal("setup failed:", err)
	}
	status := statuses[0]

	tests := statusesTests{}
	t.Run("get", tests.getStatus(status))
	t.Run("add", tests.addStatus(project))
	t.Run("del", tests.delStatus(status))
}

func (test *statusesTests) getStatus(expected Status) func(t *testing.T) {
	return func(t *testing.T) {
		actual, err := test.StatusesRepo.Get(expected.ID, expected.PID)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(expected, actual) {
			t.Error("expected/actual mismatch:", expected, actual)
		}
	}
}

func (test *statusesTests) addStatus(p Project) func(t *testing.T) {
	return func(t *testing.T) {
		s := Status{
			PID:   p.ID,
			Name:  "newstatus",
			SeqNo: 2,
		}
		id, err := test.StatusesRepo.Add(s)
		if err != nil {
			t.Fatal(err)
		}
		s.ID = id
		if !checkStatusExists(s) {
			t.Error("not found", s)
		}
	}
}

func (test *statusesTests) delStatus(s Status) func(t *testing.T) {
	return func(t *testing.T) {
		err := test.Del(s.ID, s.PID)
		if err != nil {
			t.Fatal(err)
		}
		if checkStatusExists(s) {
			t.Error("indicated ok but target entry still exists:", s)
		}
	}
}

func checkStatusExists(expected Status) bool {
	q := `SELECT * FROM statuses WHERE id=$1`
	var actual Status
	err := db.QueryRow(q, expected.ID).Scan(&actual.ID, &actual.PID, &actual.SeqNo, &actual.Name)
	return err == nil && reflect.DeepEqual(expected, actual)
}
