package middleware

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"y_finalproject/persistence"
)

func initialStatus(pid int64) persistence.TaskStatus {
	return persistence.TaskStatus{
		PID:   pid,
		Name:  "default",
		SeqNo: 0,
	}
}
func AddProject(w http.ResponseWriter, r *http.Request) {
	//todo add transaction
	var p persistence.Project
	json.NewDecoder(r.Body).Decode(&p)
	id, err := persistence.AddProject(p)
	if err == nil {
		_, err := persistence.AddTaskStatus(initialStatus(id))
		if err == nil {
			createdOk(w, id)
		}
	} else {
		reqFailed(w, err)
	}
}

func DelProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	if err := persistence.DelProject(int64(id)); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
	}
}

func ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := persistence.ListProjects()
	if err == nil {
		json.NewEncoder(w).Encode(projects)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	p, err := persistence.GetProject(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(p)
	}
}

func UpdProject(w http.ResponseWriter, r *http.Request) {
	var p persistence.Project
	json.NewDecoder(r.Body).Decode(&p)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	p.ID = int64(id)
	if persistence.UpdProject(p) != nil {
		w.WriteHeader(http.StatusNotFound)
	}
}
