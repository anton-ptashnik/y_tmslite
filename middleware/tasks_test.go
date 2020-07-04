package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"y_finalproject/persistence"
)

func TestAddTask(t *testing.T) {
	var expectedID int64 = 2
	addTaskOp := func(t persistence.Task) (int64, error) {
		return expectedID, nil
	}
	handler := AddTask(addTaskOp)
	req := httptest.NewRequest("post", "/projects/1/tasks", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	var r struct{ ID int64 }
	json.NewDecoder(resp.Body).Decode(&r)
	if resp.StatusCode != http.StatusCreated || r.ID != expectedID {
		t.Errorf("expected %v %v, but got %v %v", 201, expectedID, resp.StatusCode, r)
	}
}

func TestListTasks(t *testing.T) {
	expectedTasks := []persistence.Task{{ID: 1, ProjectID: 2, StatusID: 3, PriorityID: 1, Name: "test"}}
	listTasksFakeOp := listTasksOp(func(_ int64) ([]persistence.Task, error) {
		return expectedTasks, nil
	})

	handler := ListTasks(listTasksFakeOp)
	req := httptest.NewRequest("get", "/projects/1/tasks", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	var actualTasks []persistence.Task
	json.NewDecoder(resp.Body).Decode(&actualTasks)
	if resp.StatusCode != http.StatusOK || !reflect.DeepEqual(expectedTasks, actualTasks) {
		t.Errorf("expected %v %v, but got %v %v", 200, expectedTasks, resp.StatusCode, actualTasks)
	}
}

func TestGetTask(t *testing.T) {
	expectedTask := persistence.Task{
		ID:         2,
		ProjectID:  1,
		StatusID:   1,
		PriorityID: 0,
		Name:       "test",
	}
	getTaskFakeOp := getTaskOp(func(id int64) (persistence.Task, error) {
		//if id != expectedTask.ID {
		//	return persistence.Task{}, errors.New("")
		//}
		return expectedTask, nil
	})

	handler := GetTask(getTaskFakeOp)
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

func TestDelTask(t *testing.T) {
	delFakeOp := delTaskOp(func(_ int64) error {
		return nil
	})
	handler := DelTask(delFakeOp)
	endpoint := "/projects/1/tasks/2"
	req := httptest.NewRequest("delete", endpoint, nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Error("expected/actual status mismatch:", 200, resp.StatusCode)
	}
}

func TestUpdTask(t *testing.T) {
	updFakeOp := updTaskOp(func(_ persistence.Task) error {
		return nil
	})

	handler := UpdTask(updFakeOp)
	endpoint := "/projects/1/tasks/2"
	req := httptest.NewRequest("put", endpoint, nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Error("expected/actual status mismatch:", 200, resp.StatusCode)
	}

}
