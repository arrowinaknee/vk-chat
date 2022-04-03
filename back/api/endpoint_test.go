package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var ep = Endpoint{
	Get:    func(w http.ResponseWriter, r *http.Request) {},
	Delete: func(w http.ResponseWriter, r *http.Request) {},
}

func _testEndpoint(e Endpoint, method string) *http.Response {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, "/arbitrary", nil)

	e.ServeHTTP(w, r)

	return w.Result()
}

func TestEndpointHandled(t *testing.T) {
	res := _testEndpoint(ep, http.MethodGet)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Handled request returned %d instead of OK", res.StatusCode)
	}
}
func TestEndpointUnhandled(t *testing.T) {
	res := _testEndpoint(ep, http.MethodPost)

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("Unhandled request returned %d instead of Method Not Allowed", res.StatusCode)
	}
}
func TestEndpointUnsupported(t *testing.T) {
	res := _testEndpoint(ep, http.MethodPatch)

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("Unsupported request returned %d instead of Method Not Allowed", res.StatusCode)
	}
}
