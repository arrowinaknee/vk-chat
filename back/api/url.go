package api

import "net/http"

// shortcuts for generic get request uses
type UrlRequest struct {
	Base *http.Request
}

func (r UrlRequest) Id() string {
	return r.Base.FormValue("id")
}
func (r UrlRequest) Query() string {
	return r.Base.FormValue("q")
}

// count, skip etc.

type UrlHandler func(res *JsonResponse, r *UrlRequest)

func Url(handler UrlHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(&JsonResponse{w}, &UrlRequest{r})
	}
}
