package api

import "net/http"

type HandlerFactory[O any] func(http.HandlerFunc) O

type IdRequest[T any] struct {
	Id T `json:"id"`
}

type EmptyRequest struct{}
