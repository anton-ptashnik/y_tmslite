package middleware

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"y_finalproject/persistence"
	"y_finalproject/service"
)

type TasksService interface {
	service.TasksRepo
}

type TasksHandler struct {
	TasksService
}

func (h *TasksHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	var t persistence.Task
	json.NewDecoder(r.Body).Decode(&t)
	pid, _ := strconv.Atoi(chi.URLParam(r, "pid"))
	sid, _ := strconv.Atoi(chi.URLParam(r, "sid"))
	t.ProjectID = int64(pid)
	t.StatusID = int64(sid)
	id, err := h.TasksService.Add(t)
	if err != nil {
		reqFailed(w, err)
	} else {
		createdOk(w, id)
	}
}

func (h *TasksHandler) DelTask(w http.ResponseWriter, r *http.Request) {
	pid, _ := strconv.Atoi(chi.URLParam(r, "pid"))
	tid, _ := strconv.Atoi(chi.URLParam(r, "tid"))
	err := h.TasksService.Del(int64(tid), int64(pid))
	if err != nil {
		reqFailed(w, err)
	}
}

func (h *TasksHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	pid, _ := strconv.Atoi(chi.URLParam(r, "pid"))
	tasks, err := h.TasksService.List(int64(pid))
	if err != nil {
		reqFailed(w, err)
	} else {
		json.NewEncoder(w).Encode(tasks)
	}
}

func (h *TasksHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	pid, _ := strconv.Atoi(chi.URLParam(r, "pid"))
	tid, _ := strconv.Atoi(chi.URLParam(r, "tid"))
	task, err := h.TasksService.Get(int64(tid), int64(pid))
	if err != nil {
		reqFailed(w, err)
	} else {
		json.NewEncoder(w).Encode(task)
	}
}

func (h *TasksHandler) UpdTask(w http.ResponseWriter, r *http.Request) {
	//pid, _ := strconv.Atoi(chi.URLParam(r, "pid"))
	tid, _ := strconv.Atoi(chi.URLParam(r, "tid"))
	var t persistence.Task
	json.NewDecoder(r.Body).Decode(&t)
	t.ID = int64(tid)
	err := h.TasksService.Upd(t)
	if err != nil {
		reqFailed(w, err)
	}
}
