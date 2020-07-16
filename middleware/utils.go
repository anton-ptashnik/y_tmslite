package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

func createdOk(w http.ResponseWriter, id int64) {
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ID int64
	}{id})
}

func reqFailed(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func extractIdParam(r *http.Request, paramName string) int {
	param, err := strconv.Atoi(chi.URLParam(r, paramName))
	if err != nil {
		panic(fmt.Sprintf("%s is not provided; %v",paramName, err))
	}
	return param
}