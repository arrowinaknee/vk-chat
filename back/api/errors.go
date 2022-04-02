package api

import (
	"fmt"
	"net/http"
)

func ErrPageNotFound(w http.ResponseWriter, r *http.Response) {
	JsonResponse{w}.Error(http.StatusNotFound, "no endpoint at "+r.Request.URL.Path)
}
func ErrResNotFound(w http.ResponseWriter) {
	JsonResponse{w}.Error(http.StatusNotFound, "resource not found")
}
func ErrMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	JsonResponse{w}.Error(http.StatusMethodNotAllowed, fmt.Sprintf("method %s not allowed", r.Method))
}
