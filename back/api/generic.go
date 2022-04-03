package api

type IdRequest[T any] struct {
	Id T `json:"id"`
}

type EmptyRequest struct{}
