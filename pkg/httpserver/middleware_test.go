package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCorsHeaders(t *testing.T) {
	t.Run("returns a copy of the default CORS headers", func(t *testing.T) {
		if diff := cmp.Diff(CorsHeaders(), corsHeaders); diff != "" {
			t.Errorf("CorsHeaders() mismatch (-want +got):\n%s", diff)
		}
	})
}

func Test_headersMiddleware(t *testing.T) {
	t.Run("adds headers to the response", func(t *testing.T) {
		headers := http.Header{
			"X-Test": {"test"},
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		headersMiddleware(handler, headers).ServeHTTP(w, r)

		if diff := cmp.Diff(headers, w.Result().Header); diff != "" {
			t.Errorf("headers mismatch (-want +got):\n%s", diff)
		}
	})
}
