package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JsonResponse struct {
	w http.ResponseWriter
}

func (res JsonResponse) Write(v any) (err error) {
	err = json.NewEncoder(res.w).Encode(v)
	if err != nil {
		res.Error(http.StatusInternalServerError, "could not encode response")
	}
	return
}
func (res JsonResponse) Error(status int, err string) {
	res.w.WriteHeader(status)
	// hardcoded json generation to avoid possible loop
	fmt.Fprintf(res.w, "{\"status\":%d,\"error\":\"%s\"}", status, err)
}
func (res JsonResponse) NotFound() {
	res.Error(http.StatusNotFound, "resource not found")
}

type JsonHandler[T any] func(*JsonResponse, *T)

// Middleware that decodes json from request body before calling the handler
//
// If decoding resulted in an error, handler is not called and status
// is set to internal server error followed by a json message
func Json[T any](handler JsonHandler[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := &JsonResponse{w}
		v := new(T)

		err := json.NewDecoder(r.Body).Decode(v)

		if err != nil {
			res.Error(http.StatusBadRequest, err.Error())
			return
		}

		handler(res, v)
	}
}
