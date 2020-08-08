package middleware

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"testing"
)

type entityCreatedBody struct {
	ID int64
}

func createdOk(w http.ResponseWriter, id int64) {
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entityCreatedBody{id})
}

func reqFailed(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func clientErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

func setCtx(r *http.Request, ctx *chi.Context) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}

func checkStatus(expected int, response *httptest.ResponseRecorder, t *testing.T) {
	if response.Code != expected {
		t.Error("status code mismatch, exp/act:", expected, response.Code)
	}
}

func callHandler(handlerFunc http.HandlerFunc, r *http.Request, ctx *chi.Context) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	c := r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	handlerFunc(rec, c)
	return rec
}