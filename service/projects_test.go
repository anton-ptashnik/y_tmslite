package service

import (
	"reflect"
	"testing"
	"y_finalproject/persistence"
)

type projectTests struct {
	ProjectsService
}

func TestProjects(t *testing.T) {
	var tests projectTests

	t.Run(`add`, tests.add)
}

func (test *projectTests) add(t *testing.T) {
	expectedProj := persistence.Project{
		ID:          2,
		Name:        "newproj",
		Description: "abc",
	}
	expectedStatus := persistence.Status{
		PID:   2,
		Name:  "default",
		SeqNo: 1,
	}
	var actualProj persistence.Project
	var actualStatus persistence.Status

	test.ProjectsRepoTx = func(tx persistence.Tx) ProjectsRepo {
	 	return &fakeProjectsRepo{addOp: func(project persistence.Project) (int64, error) {
			actualProj = project
			return project.ID, nil
		}}
	}
	test.ProjectsService.InsertStatusOpTx = newFakeStatusInsertOp(&actualStatus)
	test.ProjectsService.TxInitiator = func() (persistence.Tx, error) {
		return fakeTx{}, nil
	}

	pid, err := test.ProjectsService.Add(expectedProj)
	if err != nil {
		t.Error(`project creation failed`, err)
	}
	if pid != expectedProj.ID {
		t.Error(`incorrect project ID returned`)
	}
	if !reflect.DeepEqual(expectedProj, actualProj) {
		t.Errorf(`requested/stored project mismatch; exp: %+v, act: %+v`, expectedProj, actualProj)
	}
	if !reflect.DeepEqual(expectedStatus, actualStatus) {
		t.Errorf(`wrong init status setup; exp: %+v, act: %+v`, expectedStatus, actualStatus)
	}
}

type fakeProjectsRepo struct {
	addOp func(project persistence.Project) (int64,error)
}

func (f *fakeProjectsRepo) Add(project persistence.Project) (int64, error) {
	return f.addOp(project)
}

func (f fakeProjectsRepo) Get(id int64) (persistence.Project, error) {
	panic("implement me")
}

func (f fakeProjectsRepo) Del(id int64) error {
	panic("implement me")
}

func (f fakeProjectsRepo) Upd(project persistence.Project) error {
	panic("implement me")
}

func (f fakeProjectsRepo) List() ([]persistence.Project, error) {
	panic("implement me")
}

func newFakeStatusInsertOp(s *persistence.Status) InsertStatusOpTx {
	return func(tx persistence.Tx) InsertStatusOp {
		return func(status persistence.Status) (int64, error) {
			*s = status
			return status.ID, nil
		}
	}
}