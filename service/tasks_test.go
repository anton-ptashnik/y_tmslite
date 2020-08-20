package service

import (
	"reflect"
	"testing"
	"y_finalproject/persistence"
)

func TestTasks(t *testing.T) {
	var tests tasksTests

	t.Run(`list`, tests.list)
}

type tasksTests struct {
	TasksService
}

func (test *tasksTests) list(t *testing.T) {
	tasks := []persistence.Task{
		{Name: "test", StatusID: 1},
		{Name: "target1", StatusID: 2},
		{Name: "target2", StatusID: 2},
		{Name: "test2", StatusID: 3},
		{Name: "target3", StatusID: 4},
	}
	var expected []persistence.Task
	expected = append(expected,tasks[1], tasks[2], tasks[4])
	test.TasksRepo = fakeTasksRepo{listOp: func(pid int64) ([]persistence.Task, error) {
		return tasks, nil
	}}

	res, err := test.TasksService.List(TaskFilterTemplate{
		Name:     "target",
		Statuses: []int64{2, 4},
	})

	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Error("exp/act mismatch:", expected, res)
	}
}

type fakeTasksRepo struct {
	listOp func(pid int64) ([]persistence.Task, error)
}

func (f fakeTasksRepo) List(pid int64) ([]persistence.Task, error) {
	return f.listOp(pid)
}

func (f fakeTasksRepo) Del(id int64, pid int64) error {
	panic("implement me")
}

func (f fakeTasksRepo) Add(task persistence.Task) (int64, error) {
	panic("implement me")
}

func (f fakeTasksRepo) Upd(task persistence.Task) error {
	panic("implement me")
}

func (f fakeTasksRepo) Get(id int64, pid int64) (persistence.Task, error) {
	panic("implement me")
}
