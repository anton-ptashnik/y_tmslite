package middleware

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"net/url"
	"strconv"
	"y_finalproject/persistence"
	"y_finalproject/service"
)

type TasksService interface {
	Del(id int64, pid int64) error
	Add(task persistence.Task) (int64, error)
	Upd(task persistence.Task) error
	Get(id int64, pid int64) (persistence.Task, error)
	List(filter service.TaskFilterTemplate) ([]persistence.Task, error)
}

type TasksHandler struct {
	TasksService
}

func (h *TasksHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	var t persistence.Task
	json.NewDecoder(r.Body).Decode(&t)
	pid, _ := strconv.Atoi(chi.URLParam(r, "pid"))
	t.ProjectID = int64(pid)
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
	filter := parseFilter(r.URL.Query())
	filter.Pid = int64(pid)
	tasks, err := h.TasksService.List(filter)
	if err != nil {
		reqFailed(w, err)
	} else {
		json.NewEncoder(w).Encode(tasks)
	}
}

func parseFilter(query url.Values) service.TaskFilterTemplate {
	filter := service.TaskFilterTemplate{}
	filter.Name = query.Get(`name`)
	rSids := query[`status`]
	var sids []int64
	for _, s := range rSids {
		numb, _ := strconv.Atoi(s)
		sids = append(sids, int64(numb))
	}
	filter.Statuses = sids
	return filter
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
