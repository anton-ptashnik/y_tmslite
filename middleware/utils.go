package middleware

import (
	"encoding/json"
	"net/http"
)

func createdOk(w http.ResponseWriter, id int64) {
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ID int64
	}{id})
}
