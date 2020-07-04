package persistence

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

type projectsTest struct {
	*sql.DB
}

func TestProjects(t *testing.T) {
	db := dbConn()
	defer db.Close()
	tests := projectsTest{db}
	projects, err := prepareProjects(db, 2)
	if err != nil {
		t.Fatal("preparation failed", err)
	}

	t.Run("list", tests.listProjects())
	t.Run("get", tests.getProject(projects))
	t.Run("add", tests.addProject)
	t.Run("del", tests.delProject(projects[0]))
	t.Run("upd", tests.updProject(projects[1]))
}

func prepareProjects(db *sql.DB, n int) ([]Project, error) {
	var p = Project{
		Name:        "newp",
		Description: "somethingverImportant",
	}
	q := `INSERT INTO projects (name, description) VALUES ($1,$2) RETURNING id`
	var res []Project
	for i := n; i > 0; i-- {
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

func (ps *projectsTest) addProject(t *testing.T) {
	p := Project{
		Name:        "newp",
		Description: "abc",
	}
	_, err := AddProject(p)
	if err != nil {
		t.Error(err)
	}
}

func (ps *projectsTest) getProject(projects []Project) func(t *testing.T) {
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

func (ps *projectsTest) delProject(p Project) func(*testing.T) {
	return func(t *testing.T) {
		err := DelProject(p.ID)
		if err != nil {
			t.Error(err)
		}
		if checkProjectExists(ps.DB, p) {
			t.Error("failed to remove:", p)
		}
	}
}

func (ps *projectsTest) listProjects() func(*testing.T) {
	return func(t *testing.T) {
		list, err := ListProjects()
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range list {
			if !checkProjectExists(ps.DB, v) {
				t.Error("expected entry is missing:", v)
			}
		}
	}
}

func checkProjectExists(db *sql.DB, expectedProj Project) bool {
	q := `select * from projects where id=$1`
	var actualProj Project
	err := db.QueryRow(q, expectedProj.ID).Scan(&actualProj.ID, &actualProj.Name, &actualProj.Description)
	return err == nil && reflect.DeepEqual(expectedProj, actualProj)
}

func (ps *projectsTest) updProject(p Project) func(t *testing.T) {
	upd := p
	upd.Name += "_updated"
	return func(t *testing.T) {
		err := UpdProject(upd)
		if err != nil {
			t.Error("update failed:", p)
		}
		if !checkProjectExists(ps.DB, upd) {
			t.Error("updated entry does not equal to expected:", upd)
		}
	}
}
