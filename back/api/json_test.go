package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestStruct struct {
	Text string `json:"text"`
	Num  int    `json:"num"`
}

var testVal = TestStruct{"abc", 12}

func TestJsonResponseParsable(t *testing.T) {
	w := httptest.NewRecorder()
	res := JsonResponse{w}

	in := testVal

	if err := res.Write(&in); err != nil {
		t.Fatalf("Error encoding json response: %s", err)
	}

	var out TestStruct
	err := json.NewDecoder(w.Result().Body).Decode(&out)

	if err != nil {
		t.Fatalf("Cannot parse json response: %s", err)
	}

	if out != in {
		t.Fatalf("Response did not match data")
	}
}

func TestJsonResponseErrorParsable(t *testing.T) {
	w := httptest.NewRecorder()
	res := JsonResponse{w}

	status := http.StatusTeapot
	msg := "Test error"

	res.Error(status, msg)

	var out struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}

	err := json.NewDecoder(w.Result().Body).Decode(&out)
	if err != nil {
		t.Fatalf("Cannot parse json error response: %s", err)
	}

	if out.Status != status {
		t.Fatalf("Error status did not match input")
	}
	if out.Error != msg {
		t.Fatalf("Error message did not match input")
	}
}

func TestJsonMiddlewareNormal(t *testing.T) {
	in := testVal

	buf, err := json.Marshal(&in)
	if err != nil {
		t.Fatalf("Test error: could not marshal input: %s", err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/anything", bytes.NewReader(buf))

	fun := Json(func(res *JsonResponse, r *JsonRequest[TestStruct]) {
		if *r.V != in {
			t.Fatal("Handler input did not match expected value")
		}
		res.Write(r.V)
	})
	fun(w, r)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatal("Response status was not OK")
	}
	// since JsonResponse is tested before, no need for further checks
}

func TestJsonMiddlewareBadInput(t *testing.T) {
	in := "{$arbitrary\\bad: 'input\","

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/anything", strings.NewReader(in))

	fun := Json(func(res *JsonResponse, r *JsonRequest[TestStruct]) {
		t.Fatal("Handler was called with bad request")
	})
	fun(w, r)

	res := w.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatal("Response for bad input was not Bad Request")
	}
}

func TestJsonMiddlewareBadOutput(t *testing.T) {
	in := testVal

	buf, err := json.Marshal(&in)
	if err != nil {
		t.Fatalf("Test error: could not marshal input: %s", err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/anything", bytes.NewReader(buf))

	fun := Json(func(res *JsonResponse, r *JsonRequest[TestStruct]) {
		res.Write(make(chan int))
	})
	fun(w, r)

	res := w.Result()
	// all serverside error codes match 5XX
	if res.StatusCode/100 != 5 {
		t.Fatalf("Failed response was not marked as serverside error (got %d)", res.StatusCode)
	}
}
