package middleware

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"y_finalproject/persistence"
)

func AddComment(w http.ResponseWriter, r *http.Request) {
	var c persistence.Comment
	json.NewDecoder(r.Body).Decode(&c)
	tid, _ := strconv.Atoi(chi.URLParam(r, "tid"))
	c.TaskID = int64(tid)
	id, err := persistence.AddComment(c)
	if err != nil {
		reqFailed(w, err)
	} else {
		createdOk(w, id)
	}
}

func DelComment(w http.ResponseWriter, r *http.Request) {
	cid, _ := strconv.Atoi(chi.URLParam(r, "cid"))
	err := persistence.DelComment(int64(cid))
	if err != nil {
		reqFailed(w, err)
	}
}

func ListComments(w http.ResponseWriter, r *http.Request) {
	tid, _ := strconv.Atoi(chi.URLParam(r, "tid"))
	comments, err := persistence.ListComments(int64(tid))
	if err != nil {
		reqFailed(w, err)
	} else {
		json.NewEncoder(w).Encode(comments)
	}
}

func UpdComment(w http.ResponseWriter, r *http.Request) {
	cid, _ := strconv.Atoi(chi.URLParam(r, "cid"))
	var c persistence.Comment
	json.NewDecoder(r.Body).Decode(&c)
	c.ID = int64(cid)
	err := persistence.UpdComment(c)
	if err != nil {
		reqFailed(w, err)
	}
}

func GetComment(w http.ResponseWriter, r *http.Request) {
	cid, _ := strconv.Atoi(chi.URLParam(r, "cid"))
	comment, err := persistence.GetComment(int64(cid))
	if err != nil {
		reqFailed(w, err)
	} else {
		json.NewEncoder(w).Encode(comment)
	}
}
