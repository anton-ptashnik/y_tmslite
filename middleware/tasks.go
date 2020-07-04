package middleware

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"y_finalproject/persistence"
)

type addTaskOp func(p persistence.Task) (int64, error)
type listTasksOp func(pid int64) ([]persistence.Task, error)
type getTaskOp func(id int64) (persistence.Task, error)
type delTaskOp func(id int64) error
type updTaskOp func(t persistence.Task) error

func AddTask(op addTaskOp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t persistence.Task
		json.NewDecoder(r.Body).Decode(&t)
		pid, _ := strconv.Atoi(chi.URLParam(r, "pid"))
		t.ProjectID = int64(pid)
		id, err := op(t)
		if err != nil {
			reqFailed(w, err)
		} else {
			createdOk(w, id)
		}
	}
}

func DelTask(op delTaskOp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//pid, _ := strconv.Atoi(chi.URLParam(r, "pid"))
		tid, _ := strconv.Atoi(chi.URLParam(r, "tid"))
		err := op(int64(tid))
		if err != nil {
			reqFailed(w, err)
		}
	}
}

func ListTasks(op listTasksOp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pid, _ := strconv.Atoi(chi.URLParam(r, "pid"))
		tasks, err := op(int64(pid))
		if err != nil {
			reqFailed(w, err)
		} else {
			json.NewEncoder(w).Encode(tasks)
		}
	}
}

func GetTask(op getTaskOp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//pid, _ := strconv.Atoi(chi.URLParam(r, "pid"))
		tid, _ := strconv.Atoi(chi.URLParam(r, "tid"))
		task, err := op(int64(tid))
		if err != nil {
			reqFailed(w, err)
		} else {
			json.NewEncoder(w).Encode(task)
		}
	}
}

func UpdTask(op updTaskOp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//pid, _ := strconv.Atoi(chi.URLParam(r, "pid"))
		tid, _ := strconv.Atoi(chi.URLParam(r, "tid"))
		var t persistence.Task
		json.NewDecoder(r.Body).Decode(&t)
		t.ID = int64(tid)
		err := op(t)
		if err != nil {
			reqFailed(w, err)
		}
	}
}
