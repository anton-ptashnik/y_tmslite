package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type statusesTests struct {
	*sql.DB
}

func TestStatuses(t *testing.T) {
	db := dbConn()
	defer db.Close()
	p, err := prepareStatuses(db)
	if err != nil {
		t.Fatal("setup failed:", err)
	}
	tests := statusesTests{db}

	t.Run("get", tests.getStatus(p.TaskStatuses[0]))
	t.Run("add", tests.addStatus(p))
	t.Run("del", tests.delStatus(p.TaskStatuses[0]))
}

func prepareStatuses(db *sql.DB) (Project, error) {
	projects, err := prepareProjects(db, 1)
	if err != nil {
		return Project{}, errors.New(fmt.Sprint("projects prep failed:", err))
	}
	p := projects[0]
	var status = TaskStatus{
		PID:   p.ID,
		Name:  "default",
		SeqNo: 0,
	}
	q := `INSERT INTO statuses (pid, name, seqNo) VALUES ($1,$2,$3) RETURNING id`
	err = db.QueryRow(q, status.PID, status.Name, status.SeqNo).Scan(&status.ID)
	if err != nil {
		return Project{}, err
	}
	p.TaskStatuses = append(p.TaskStatuses, status)
	return p, nil
}

func (st *statusesTests) getStatus(expected TaskStatus) func(t *testing.T) {
	return func(t *testing.T) {
		actual, err := GetTaskStatus(expected.ID)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(expected, actual) {
			t.Error("expected/actual mismatch:", expected, actual)
		}
	}
}

func (st *statusesTests) addStatus(p Project) func(t *testing.T) {
	return func(t *testing.T) {
		s := TaskStatus{
			PID:   p.ID,
			Name:  "newstatus",
			SeqNo: 2,
		}
		id, err := AddTaskStatus(s)
		if err != nil {
			t.Fatal(err)
		}
		s.ID = id
		if !checkStatusExists(s) {
			t.Error("not found", s)
		}
	}
}

func (st *statusesTests) delStatus(s TaskStatus) func(t *testing.T) {
	return func(t *testing.T) {
		err := DelTaskStatus(s.ID)
		if err != nil {
			t.Fatal(err)
		}
		if checkStatusExists(s) {
			t.Error("indicated ok but target entry still exists:", s)
		}
	}
}

func checkStatusExists(expected TaskStatus) bool {
	db := dbConn()
	defer db.Close()
	q := `SELECT * FROM statuses WHERE id=$1`
	var actual TaskStatus
	err := db.QueryRow(q, expected.ID).Scan(&actual.ID, &actual.PID, &actual.SeqNo, &actual.Name)
	return err == nil && reflect.DeepEqual(expected, actual)
}
