package api

import (
	"net/http"
)

type Endpoint struct {
	Get    http.HandlerFunc
	Post   http.HandlerFunc
	Delete http.HandlerFunc
	Put    http.HandlerFunc
}

func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if e.Get != nil {
			e.Get(w, r)
		} else {
			ErrMethodNotAllowed(w, r)
		}
	case http.MethodPost:
		if e.Post != nil {
			e.Post(w, r)
		} else {
			ErrMethodNotAllowed(w, r)
		}
	case http.MethodDelete:
		if e.Delete != nil {
			e.Delete(w, r)
		} else {
			ErrMethodNotAllowed(w, r)
		}
	case http.MethodPut:
		if e.Put != nil {
			e.Put(w, r)
		} else {
			ErrMethodNotAllowed(w, r)
		}
	default:
		ErrMethodNotAllowed(w, r)
	}
}
