package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/extreme-business/lingo/pkg/httpserver"
	"github.com/google/go-cmp/cmp"
)

func TestCorsHeaders(t *testing.T) {
	t.Run("returns a copy of the default CORS headers", func(t *testing.T) {
		corsHeaders := http.Header{
			"Access-Control-Allow-Origin": {"*"},
			"Access-Control-Allow-Methods": {
				strings.Join([]string{
					http.MethodGet,
					http.MethodPost,
					http.MethodPut,
					http.MethodDelete,
				}, ",")},
			"Access-Control-Allow-Headers": {"Content-Type, Accountorization"},
		}

		if diff := cmp.Diff(httpserver.CorsHeaders(), corsHeaders); diff != "" {
			t.Errorf("CorsHeaders() mismatch (-want +got):\n%s", diff)
		}
	})
}

func Test_headersMiddleware(t *testing.T) {
	t.Run("adds headers to the response", func(t *testing.T) {
		headers := http.Header{
			"X-Test": {"test"},
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		httpserver.HeadersMiddleware(handler, headers).ServeHTTP(w, r)

		if diff := cmp.Diff(headers, w.Result().Header); diff != "" {
			t.Errorf("headers mismatch (-want +got):\n%s", diff)
		}
	})
}
