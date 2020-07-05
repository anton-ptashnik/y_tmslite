package middleware

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"y_finalproject/persistence"
)

var (
	errLastStatus = errors.New("last status cannot be deleted")
)

func AddTaskStatus(w http.ResponseWriter, r *http.Request) {
	pid, _ := strconv.Atoi(chi.URLParam(r, "pid"))
	var d persistence.TaskStatus
	json.NewDecoder(r.Body).Decode(&d)
	d.PID = int64(pid)
	id, err := persistence.AddTaskStatus(d)
	if err != nil {
		reqFailed(w, err)
	} else {
		createdOk(w, id)
	}

}
func DelTaskStatus(w http.ResponseWriter, r *http.Request) {
	// todo introduce service layer delete operation
	pid, _ := strconv.Atoi(chi.URLParam(r, "pid"))
	sid, _ := strconv.Atoi(chi.URLParam(r, "sid"))
	statuses, err := persistence.ListStatuses(int64(pid))
	if err != nil {
		reqFailed(w, err)
		return
	}
	if len(statuses) == 1 && statuses[0].ID == int64(sid) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errLastStatus.Error()))
		return
	}
	if err := persistence.DelTaskStatus(int64(sid)); err != nil {
		reqFailed(w, err)
	}
}

func UpdTaskStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "sid"))
	var s persistence.TaskStatus
	json.NewDecoder(r.Body).Decode(&s)
	s.ID = int64(id)
	if err := persistence.UpdStatus(s); err != nil {
		reqFailed(w, err)
	}
}

func ListTaskStatuses(w http.ResponseWriter, r *http.Request) {
	pid, _ := strconv.Atoi(chi.URLParam(r, "pid"))
	res, err := persistence.ListStatuses(int64(pid))
	if err != nil {
		reqFailed(w, err)
	} else {
		json.NewEncoder(w).Encode(res)
	}
}

func GetStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "sid"))
	status, err := persistence.GetTaskStatus(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(status)
	}
}
