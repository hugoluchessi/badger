package badger

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildHandler(t *testing.T) {
	request, _ := http.NewRequest("GET", "/health-check", nil)
	response := httptest.NewRecorder()

	headerkey := "X-Middleware"
	headerkey2 := "X-Middleware-2"
	expectedheadervalue := "YEAH!"

	mw := &middleware{func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Header().Add(headerkey, expectedheadervalue)
			h.ServeHTTP(res, req)
			res.Header().Add(headerkey2, expectedheadervalue)
		})
	}}

	h := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		headervalue := res.Header().Get(headerkey)
		headervalue2 := res.Header().Get(headerkey2)

		if headervalue != headervalue {
			t.Errorf("Test failed, invalid '%s' header value, got '%s' expected '%s'.", headerkey, headervalue, expectedheadervalue)
		}

		if headervalue2 != "" {
			t.Errorf("Test failed, invalid '%s' header value, got '%s' expected ''.", headerkey2, headervalue2)
		}
	})

	builthandler := mw.BuildHandler(h)

	builthandler.ServeHTTP(response, request)
}
