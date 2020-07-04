package persistence

import (
	"database/sql"
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

	tests := statusesTests{db}
	t.Run("get", tests.getStatus(status))
	t.Run("add", tests.addStatus(project))
	t.Run("del", tests.delStatus(status))
}

func prepareStatuses(db *sql.DB, p Project, n int) ([]Status, error) {
	var baseStatus = Status{
		PID:   p.ID,
		Name:  "default",
		SeqNo: 0,
	}
	q := `INSERT INTO statuses (pid, name, seqNo) VALUES ($1,$2,$3) RETURNING id`
	var statuses []Status
	for n > 0 {
		n--
		status := baseStatus
		status.Name = fmt.Sprint(status.Name, n)
		status.SeqNo = n
		err := db.QueryRow(q, status.PID, status.Name, status.SeqNo).Scan(&status.ID)
		if err != nil {
			return statuses, err
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}

func (st *statusesTests) getStatus(expected Status) func(t *testing.T) {
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
		s := Status{
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

func (st *statusesTests) delStatus(s Status) func(t *testing.T) {
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

func checkStatusExists(expected Status) bool {
	db := dbConn()
	defer db.Close()
	q := `SELECT * FROM statuses WHERE id=$1`
	var actual Status
	err := db.QueryRow(q, expected.ID).Scan(&actual.ID, &actual.PID, &actual.SeqNo, &actual.Name)
	return err == nil && reflect.DeepEqual(expected, actual)
}
