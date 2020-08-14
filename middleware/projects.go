package middleware

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"y_finalproject/persistence"
	"y_finalproject/service"
)

type ProjectsService interface {
	service.ProjectsRepo
}

type ProjectsHandler struct {
	ProjectsService
}

func (h *ProjectsHandler) AddProject(w http.ResponseWriter, r *http.Request) {
	var p persistence.Project
	json.NewDecoder(r.Body).Decode(&p)
	id, err := h.Add(p)
	if err == nil {
		createdOk(w, id)
	} else {
		reqFailed(w, err)
	}
}
func (h *ProjectsHandler) DelProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	if err := h.Del(int64(id)); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
	}
}

func (h *ProjectsHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.List()
	if err == nil {
		json.NewEncoder(w).Encode(projects)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *ProjectsHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	p, err := h.Get(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(p)
	}
}

func (h *ProjectsHandler) UpdProject(w http.ResponseWriter, r *http.Request) {
	var p persistence.Project
	json.NewDecoder(r.Body).Decode(&p)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	p.ID = int64(id)
	if h.Upd(p) != nil {
		w.WriteHeader(http.StatusNotFound)
	}
}
