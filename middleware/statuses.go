package middleware

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"y_finalproject/persistence"
	"y_finalproject/service"
)

var (
	errLastStatus = errors.New("last status cannot be deleted")
)

type StatusesHandler struct {
	service.StatusesService
}

func (h *StatusesHandler) AddStatus(w http.ResponseWriter, r *http.Request) {
	pid, _ := strconv.Atoi(chi.URLParam(r, "pid"))
	var d persistence.Status
	json.NewDecoder(r.Body).Decode(&d)
	d.PID = int64(pid)
	id, err := h.StatusesService.Add(d)
	if err != nil {
		reqFailed(w, err)
	} else {
		createdOk(w, id)
	}

}
func (h *StatusesHandler) DelStatus(w http.ResponseWriter, r *http.Request) {
	pid, _ := strconv.Atoi(chi.URLParam(r, "pid"))
	sid, _ := strconv.Atoi(chi.URLParam(r, "sid"))
	if err := h.StatusesService.Del(int64(sid), int64(pid)); err != nil {
		reqFailed(w, err)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *StatusesHandler) UpdStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "sid"))
	var s persistence.Status
	json.NewDecoder(r.Body).Decode(&s)
	s.ID = int64(id)
	if err := h.StatusesService.Upd(s); err != nil {
		reqFailed(w, err)
	}
}

func (h *StatusesHandler) ListTaskStatuses(w http.ResponseWriter, r *http.Request) {
	pid, _ := strconv.Atoi(chi.URLParam(r, "pid"))
	res, err := h.StatusesService.List(int64(pid))
	if err != nil {
		reqFailed(w, err)
	} else {
		json.NewEncoder(w).Encode(res)
	}
}

func (h *StatusesHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "sid"))
	status, err := h.StatusesService.Get(int64(id), 0)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(status)
	}
}
