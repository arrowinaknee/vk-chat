package api

import "net/http"

// shortcuts for generic get request uses
type GetRequest struct {
	Raw *http.Request
}

func (r GetRequest) Id() string {
	return r.Raw.FormValue("id")
}
func (r GetRequest) Query() string {
	return r.Raw.FormValue("id")
}

// count, skip etc.

type GetHandler func(res *JsonResponse, r *GetRequest)

func Get(handler GetHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(&JsonResponse{w}, &GetRequest{r})
	}
}
