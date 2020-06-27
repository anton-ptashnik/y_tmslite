package persistence

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

var db *sql.DB

func TestProjects(t *testing.T) {
	db = dbConn()
	defer db.Close()
	entries, err := dbSetup()
	if err != nil {
		t.Fatal("preparation failed", err)
	}

	t.Run("list", testListProjects())
	t.Run("get", testGetProject(entries))
	t.Run("add", testAddProject)
	t.Run("del", testDelProject(entries[0]))
	t.Run("upd", testUpdProject(entries[1]))
}

func dbSetup() ([]Project, error) {
	var p = Project{
		Name:        "newp",
		Description: "somethingverImportant",
	}
	q := `INSERT INTO projects (name, description) VALUES ($1,$2) RETURNING id`
	var res []Project
	for i := 2; i > 0; i-- {
		pc := p
		pc.Name += fmt.Sprint(i)
		err := db.QueryRow(q, pc.Name, pc.Description).Scan(&pc.ID)
		if err != nil {
			return nil, err
		}
		res = append(res, pc)
	}
	return res, nil
}

func testAddProject(t *testing.T) {
	p := Project{
		Name:        "newp",
		Description: "abc",
	}
	_, err := AddProject(p)
	if err != nil {
		t.Error(err)
	}
}

func testGetProject(projects []Project) func(t *testing.T) {
	return func(t *testing.T) {
		for _, expectedProj := range projects {
			actualProj, err := GetProject(expectedProj.ID)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(expectedProj, actualProj) {
				t.Error("expected/actual entry mismatch:", expectedProj, actualProj)
			}
		}
	}
}

func testDelProject(p Project) func(*testing.T) {
	return func(t *testing.T) {
		err := DelProject(p.ID)
		if err != nil {
			t.Error(err)
		}
		if dbCheckEntryExists(p) {
			t.Error("failed to remove:", p)
		}
	}
}

func testListProjects() func(*testing.T) {
	return func(t *testing.T) {
		list, err := ListProjects()
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range list {
			if !dbCheckEntryExists(v) {
				t.Error("expected entry is missing:", v)
			}
		}
	}
}

func dbCheckEntryExists(expectedProj Project) bool {
	q := `select * from projects where id=$1`
	var actualProj Project
	err := db.QueryRow(q, expectedProj.ID).Scan(&actualProj.ID, &actualProj.Name, &actualProj.Description)
	return err == nil && reflect.DeepEqual(expectedProj, actualProj)
}

func testUpdProject(p Project) func(t *testing.T) {
	upd := p
	upd.Name += "_updated"
	return func(t *testing.T) {
		err := UpdProject(upd)
		if err != nil {
			t.Error("update failed:", p)
		}
		if !dbCheckEntryExists(upd) {
			t.Error("updated entry does not equal to expected:", upd)
		}
	}
}
