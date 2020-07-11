package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"y_finalproject/persistence"
)

type tasksHandlerTest struct {
	TasksHandler
}
type fakeTasksService struct {
	tasks []persistence.Task
}

func TestTasks(t *testing.T) {
	tests := tasksHandlerTest{}

	t.Run("add", tests.addTask)
	t.Run("del", tests.delTask)
	t.Run("list", tests.listTasks)
	t.Run("upd", tests.updTask)
	t.Run("get", tests.getTask)
}
func (f *fakeTasksService) List(pid int64) ([]persistence.Task, error) {
	return f.tasks, nil
}

func (f *fakeTasksService) Del(id int64, pid int64) error {
	return nil
}

func (f *fakeTasksService) Add(task persistence.Task) (int64, error) {
	id := len(f.tasks)
	f.tasks = append(f.tasks, task)
	return int64(id), nil
}

func (f *fakeTasksService) Upd(task persistence.Task) error {
	if task.ID >= int64(len(f.tasks)) {
		return errors.New("wrong ID")
	}
	f.tasks[task.ID] = task
	return nil
}

func (f *fakeTasksService) Get(id int64, pid int64) (persistence.Task, error) {
	if id >= int64(len(f.tasks)) {
		return persistence.Task{}, errors.New("wrong ID")
	}
	return f.tasks[id], nil
}

func (test *tasksHandlerTest) addTask(t *testing.T) {
	task := persistence.Task{
		ProjectID:   0,
		StatusID:    0,
		PriorityID:  0,
		Name:        "abc",
		Description: "",
	}
	test.initTaskService([]persistence.Task{})
	req := httptest.NewRequest("post", "/projects/1/tasks", nil)
	w := httptest.NewRecorder()
	test.TasksHandler.AddTask(w, req)

	resp := w.Result()
	var r struct{ ID int64 }
	json.NewDecoder(resp.Body).Decode(&r)
	if resp.StatusCode != http.StatusCreated || r.ID != task.ID {
		t.Errorf("expected %v %v, but got %v %v", 201, task.ID, resp.StatusCode, r.ID)
	}
}

func (test *tasksHandlerTest) listTasks(t *testing.T) {
	expectedTasks := []persistence.Task{{ProjectID: 2, StatusID: 3, PriorityID: 1, Name: "test"}}

	expectedTasks = test.initTaskService(expectedTasks)
	req := httptest.NewRequest("get", "/projects/1/tasks", nil)
	w := httptest.NewRecorder()
	test.TasksHandler.ListTasks(w, req)

	resp := w.Result()
	var actualTasks []persistence.Task
	json.NewDecoder(resp.Body).Decode(&actualTasks)
	if resp.StatusCode != http.StatusOK || !reflect.DeepEqual(expectedTasks, actualTasks) {
		t.Errorf("expected %v %v, but got %v %v", 200, expectedTasks, resp.StatusCode, actualTasks)
	}
}

func (test *tasksHandlerTest) getTask(t *testing.T) {
	expectedTask := persistence.Task{
		ProjectID:  1,
		StatusID:   1,
		PriorityID: 0,
		Name:       "test",
	}
	expectedTask = test.initTaskService([]persistence.Task{
		expectedTask,
	})[0]
	handler := test.TasksHandler.GetTask

	endpoint := "/projects/1/tasks/2"
	req := httptest.NewRequest("get", endpoint, nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	var actual persistence.Task
	json.NewDecoder(resp.Body).Decode(&actual)
	if resp.StatusCode != http.StatusOK || !reflect.DeepEqual(expectedTask, actual) {
		t.Errorf("expected %v %v, but got %v %v", 200, expectedTask, resp.StatusCode, actual)
	}
}

func (test *tasksHandlerTest) delTask(t *testing.T) {
	task := test.initTaskService([]persistence.Task{
		{
			Name: "testtask",
		},
	})[0]

	handler := test.TasksHandler.DelTask
	endpoint := fmt.Sprint("/projects/0/tasks/",task.ID)
	req := httptest.NewRequest("delete", endpoint, nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Error("expected/actual status mismatch:", 200, resp.StatusCode)
	}
}

func (test *tasksHandlerTest) updTask(t *testing.T) {
	handler := test.TasksHandler.UpdTask
	endpoint := "/projects/1/tasks/2"
	req := httptest.NewRequest("put", endpoint, nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Error("expected/actual status mismatch:", 200, resp.StatusCode)
	}

}

func (test *tasksHandlerTest) initTaskService(tasks []persistence.Task) []persistence.Task {
	taskService := fakeTasksService{}
	for i := range tasks {
		tasks[i].ID, _ = taskService.Add(tasks[i])
	}
	test.TasksService = &taskService
	return tasks
}