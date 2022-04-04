package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JsonResponse struct {
	W http.ResponseWriter
}

func (res JsonResponse) Write(v any) (err error) {
	res.W.Header().Add("content-type", "application/json; charset=UTF-8")
	err = json.NewEncoder(res.W).Encode(v)
	if err != nil {
		res.Error(http.StatusInternalServerError, "could not encode response")
	}
	return
}
func (res JsonResponse) Error(status int, err string) {
	res.W.Header().Add("content-type", "application/json; charset=UTF-8")
	res.W.WriteHeader(status)
	// hardcoded json generation to avoid possible loop
	fmt.Fprintf(res.W, "{\"status\":%d,\"error\":\"%s\"}", status, err)
}
func (res JsonResponse) NotFound() {
	res.Error(http.StatusNotFound, "resource not found")
}

type JsonRequest[T any] struct {
	Base *http.Request
	V    *T // pre-parsed request body
}

type JsonHandler[T any] func(*JsonResponse, *JsonRequest[T])

// Middleware that decodes json from request body before calling the handler
//
// If decoding resulted in an error, handler is not called and status
// is set to internal server error followed by a json message
func Json[T any](handler JsonHandler[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := &JsonResponse{w}
		req := &JsonRequest[T]{r, new(T)}

		err := json.NewDecoder(r.Body).Decode(req.V)

		if err != nil {
			res.Error(http.StatusBadRequest, err.Error())
			return
		}

		handler(res, req)
	}
}
