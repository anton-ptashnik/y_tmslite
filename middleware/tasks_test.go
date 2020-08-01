package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"y_finalproject/persistence"
)

type tasksHandlerTest struct {
	TasksHandler
}

type fakeTasksService struct {
	addOp func(task persistence.Task) (int64, error)
	listOp func(pid int64) ([]persistence.Task,error)
	delOp func(id int64, pid int64) error
	updOp func(task persistence.Task) error
	getOp func(id int64, pid int64) (persistence.Task,error)
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
	return f.listOp(pid)
}

func (f *fakeTasksService) Del(id int64, pid int64) error {
	return f.delOp(id, pid)
}

func (f *fakeTasksService) Add(task persistence.Task) (int64, error) {
	return f.addOp(task)
}

func (f *fakeTasksService) Upd(task persistence.Task) error {
	return f.updOp(task)
}

func (f *fakeTasksService) Get(id int64, pid int64) (persistence.Task, error) {
	return f.getOp(id, pid)
}

func (test *tasksHandlerTest) addTask(t *testing.T) {
	body := `{	Name: 'abc', ProjectID: 1, StatusID: 0, 'PriorityID': 0 }`
	var newTaskObj persistence.Task
	json.NewDecoder(strings.NewReader(body)).Decode(&newTaskObj)

	var tasks []persistence.Task
	addTaskOp := func(task persistence.Task) (int64,error) {
		tasks = append(tasks, task)
		return task.ID,nil
	}
	test.TasksService = &fakeTasksService{
		addOp: addTaskOp,
	}

	// /projects/1/tasks
	ctx := context.WithValue(context.Background(), `pid`, 1)
	req, _ := http.NewRequestWithContext(ctx, "post", "/", strings.NewReader(body))
	w := httptest.NewRecorder()
	test.TasksHandler.AddTask(w, req)

	resp := w.Result()
	var r struct{ ID int64 }
	json.NewDecoder(resp.Body).Decode(&r)
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected code %v, but got %v", 201, resp.StatusCode)
	}
	if len(tasks) == 0 || !reflect.DeepEqual(tasks[0], newTaskObj) {
		t.Errorf("repo expected to contain %v, but not found. Repo: %v", newTaskObj, tasks)
	}
}

func (test *tasksHandlerTest) listTasks(t *testing.T) {
	var pid int64 = 2
	expectedTasks := []persistence.Task{{ProjectID: pid, StatusID: 3, PriorityID: 1, Name: "test"}}

	listTasksOp := func(_pid int64) ([]persistence.Task,error) {
		if _pid != pid {
			return nil, errors.New("not found")
		}
		return expectedTasks, nil
	}

	test.TasksService = &fakeTasksService{
		listOp: listTasksOp,
	}
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("pid", strconv.FormatInt(pid, 10))
	req := httptest.NewRequest("get", "/", nil)
	w := httptest.NewRecorder()
	test.TasksHandler.ListTasks(w, setCtx(req, ctx))

	resp := w.Result()
	var actualTasks []persistence.Task
	json.NewDecoder(resp.Body).Decode(&actualTasks)
	if resp.StatusCode != http.StatusOK || !reflect.DeepEqual(expectedTasks, actualTasks) {
		t.Errorf("expected %v %v, but got %v %v", 200, expectedTasks, resp.StatusCode, actualTasks)
	}
}

func (test *tasksHandlerTest) getTask(t *testing.T) {
	expectedTask := persistence.Task{
		ID: 4,
		ProjectID:  1,
		Name:       "test",
	}
	getTaskOp := func(_tid int64, _pid int64) (persistence.Task,error) {
		if _tid != expectedTask.ID || _pid != expectedTask.ProjectID {
			return persistence.Task{}, errors.New("not found")
		}
		return expectedTask, nil
	}
	test.TasksService = &fakeTasksService{
		getOp: getTaskOp,
	}
	handler := test.TasksHandler.GetTask

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("tid", strconv.FormatInt(expectedTask.ID, 10))
	ctx.URLParams.Add("pid", strconv.FormatInt(expectedTask.ProjectID, 10))
	req := httptest.NewRequest("get", "/", nil)
	w := httptest.NewRecorder()
	handler(w, setCtx(req, ctx))

	resp := w.Result()
	var actual persistence.Task
	json.NewDecoder(resp.Body).Decode(&actual)
	if resp.StatusCode != http.StatusOK || !reflect.DeepEqual(expectedTask, actual) {
		t.Errorf("expected %v %v, but got %v %v", 200, expectedTask, resp.StatusCode, actual)
	}
}

func (test *tasksHandlerTest) delTask(t *testing.T) {
	var tid int64 = 4
	var pid int64 = 5
	deleted := false
	delTaskOp := func(_tid int64, _pid int64) error {
		if _tid != tid || _pid != pid {
			return errors.New("not found")
		}
		deleted = true
		return nil
	}
	test.TasksService = &fakeTasksService{
		delOp: delTaskOp,
	}

	handler := test.TasksHandler.DelTask

	req := httptest.NewRequest("delete", "/", nil)
	w := httptest.NewRecorder()

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("pid", strconv.FormatInt(pid, 10))
	ctx.URLParams.Add("tid", strconv.FormatInt(tid, 10))
	handler(w, setCtx(req, ctx))

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Error("expected/actual status mismatch:", 200, resp.StatusCode)
	}
	if !deleted {
		t.Error("indicated ok but entry still exists")
	}
}

func (test *tasksHandlerTest) updTask(t *testing.T) {
	task := persistence.Task{
		StatusID:   1,
		PriorityID: 7,
		Name:       "test",
	}
	updTaskOp := func(_task persistence.Task) error {
		if _task.ID != task.ID || _task.ProjectID != task.ProjectID {
			return errors.New("not found")
		}
		task = _task
		return nil
	}
	getTaskOp := func(tid int64, pid int64) (persistence.Task,error) {
		return task, nil
	}
	test.TasksService = &fakeTasksService{
		updOp: updTaskOp,
		getOp: getTaskOp,
	}
	body := `{ "Name": "newName" }`
	expectedTask := persistence.Task{
		Name: "newName",
	}
	handler := test.TasksHandler.UpdTask

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add(`pid`, strconv.FormatInt(task.ProjectID, 10))
	ctx.URLParams.Add(`tid`, strconv.FormatInt(task.ID, 10))
	req := httptest.NewRequest("put", `/`, strings.NewReader(body))
	w := httptest.NewRecorder()
	handler(w, setCtx(req, ctx))

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Error("expected/actual status mismatch:", 200, resp.StatusCode)
	}
	if !reflect.DeepEqual(expectedTask, task) {
		t.Error("indicated ok but seem was not updated, expected/actual:", expectedTask, task)
	}
}

func setCtx(r *http.Request, ctx *chi.Context) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}