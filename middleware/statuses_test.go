package middleware

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"y_finalproject/persistence"
)

type statusesTests struct {
	StatusesHandler
}

func TestStatuses(t *testing.T) {
	tests := statusesTests{}

	t.Run(`add`, tests.add)
	t.Run(`del`, tests.del)
	t.Run(`del last status is prohibited`, tests.delLastStatus)
	t.Run(`upd`, tests.upd)
	t.Run(`get`, tests.get)
	t.Run(`list`, tests.list)
}

func (test *statusesTests) add(t *testing.T) {
	handler := test.StatusesHandler.AddStatus
	body := strings.NewReader(`{"Name": "newStatus", "SeqNo": 1}`)
	expectedStatus := persistence.Status{
		Name:  "newStatus",
		SeqNo: 1,
		PID: 1,
	}

	var expected = struct {
		code int
		entityID int64
	}{http.StatusCreated, 5}
	var actualStatus persistence.Status
	service := fakeStatusesService{addOp: func(status persistence.Status) (int64, error) {
		actualStatus = status
		return expected.entityID, nil
	}}
	test.StatusesService = service

	req := httptest.NewRequest("post", "/", body)
	rRec := httptest.NewRecorder()
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add(`pid`, `1`)
	handler(rRec, setCtx(req, ctx))

	checkStatus(expected.code, rRec, t)
	var rBody entityCreatedBody
	json.NewDecoder(rRec.Body).Decode(&rBody)
	if rBody.ID != expected.entityID {
		t.Error("new entity ID mismatch")
	}
	if actualStatus != expectedStatus {
		t.Errorf(`saved/requested status mismatch: %+v, %+v`, expectedStatus,actualStatus)
	}
}

func (test *statusesTests) del(t *testing.T) {
	handler := test.StatusesHandler.DelStatus
	targetStatus := persistence.Status{
		ID:  1,
		PID: 2,
	}
	var deletedStatus *persistence.Status
	test.StatusesService = fakeStatusesService{
		delOp: func(id, pid int64) error {
			deletedStatus = &persistence.Status{
				ID:  id,
				PID: pid,
			}
			return nil
		},
	}

	req := httptest.NewRequest("post", "/", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add(`sid`, strconv.FormatInt(targetStatus.ID, 10))
	ctx.URLParams.Add(`pid`, strconv.FormatInt(targetStatus.PID, 10))
	rRec := callHandler(handler, req, ctx)

	checkStatus(200, rRec, t)

	if !reflect.DeepEqual(*deletedStatus, targetStatus) {
		t.Error(`requested entry was not deleted`)
	}
}

func (test *statusesTests) upd(t *testing.T) {
	actual := persistence.Status{
		ID:    1,
		PID:   2,
		Name:  "teststatus",
		SeqNo: 1,
	}
	body := `{"Name": "upd_teststatus", "SeqNo": 1}`
	updated := persistence.Status{
		ID:    1,
		PID:   2,
		Name:  "upd_teststatus",
		SeqNo: 1,
	}
	test.StatusesService = fakeStatusesService{
		updOp: func(status persistence.Status) error {
			actual = status
			return nil
		},
	}

	req := httptest.NewRequest(`put`, `/`, strings.NewReader(body))
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add(`pid`, strconv.FormatInt(actual.PID, 10))
	ctx.URLParams.Add(`sid`, strconv.FormatInt(actual.ID, 10))
	resp := callHandler(test.StatusesHandler.UpdStatus, req, ctx)

	checkStatus(200, resp, t)
	if actual != updated {
		t.Errorf(`unexpected post upd state; exp/act: %+v / %+v`, updated, actual)
	}
}

func (test *statusesTests) get(t *testing.T) {
	expected := persistence.Status{
		ID:    1,
		PID:   2,
		Name:  "teststatus",
		SeqNo: 2,
	}

	test.StatusesService = fakeStatusesService{
		getOp: func(id, pid int64) (persistence.Status, error) {
			return expected, nil
		},
	}

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add(`id`, strconv.FormatInt(expected.ID, 10))
	ctx.URLParams.Add(`pid`, strconv.FormatInt(expected.PID, 10))
	req := httptest.NewRequest(`get`, `/`, nil)
	resp := callHandler(test.GetStatus, req, ctx)

	checkStatus(200, resp, t)
	var actual persistence.Status
	json.NewDecoder(resp.Body).Decode(&actual)

	if actual != expected {
		t.Errorf("exp/act: %+v / %+v", expected, actual)
	}
}

func (test *statusesTests) delLastStatus(t *testing.T) {
	test.StatusesService = fakeStatusesService{
		delOp: func(id, pid int64) error {
			return errLastStatus
		},
	}

	req := httptest.NewRequest("post", "/", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add(`sid`, strconv.FormatInt(1, 10))
	ctx.URLParams.Add(`pid`, strconv.FormatInt(2, 10))
	rRec := callHandler(test.StatusesHandler.DelStatus, req, ctx)

	checkStatus(400, rRec, t)
}

func (test *statusesTests) list(t *testing.T) {
	expected := []persistence.Status{
		{
			ID:   1,
			PID:  2,
			Name: "s1",
		},
		{
			ID:    2,
			PID:   2,
			Name:  "s2",
		},
	}
	test.StatusesService = fakeStatusesService{
		listOp: func(pid int64) ([]persistence.Status, error) {
			return expected, nil
		},
	}
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add(`pid`, strconv.Itoa(2))
	req := httptest.NewRequest(`get`, `/`, nil)
	resp := callHandler(test.ListStatuses, req, ctx)

	checkStatus(200, resp, t)
	var actual []persistence.Status
	json.NewDecoder(resp.Body).Decode(&actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Error(`exp/act:`, expected, actual)
	}
}

type fakeStatusesService struct {
	addOp func(status persistence.Status) (int64,error)
	delOp func(id, pid int64) error
	updOp func(status persistence.Status) error
	getOp func(id, pid int64) (persistence.Status, error)
	listOp func(pid int64) ([]persistence.Status, error)
}

func (f fakeStatusesService) List(pid int64) ([]persistence.Status, error) {
	return f.listOp(pid)
}

func (f fakeStatusesService) Del(id int64, pid int64) error {
	return f.delOp(id, pid)
}

func (f fakeStatusesService) Add(status persistence.Status) (int64, error) {
	return f.addOp(status)
}

func (f fakeStatusesService) Upd(status persistence.Status) error {
	return f.updOp(status)
}

func (f fakeStatusesService) Get(id int64, pid int64) (persistence.Status, error) {
	return f.getOp(id, pid)
}
