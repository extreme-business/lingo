package httpmiddleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/extreme-business/lingo/pkg/httpmiddleware"
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

		if diff := cmp.Diff(httpmiddleware.CorsHeaders(), corsHeaders); diff != "" {
			t.Errorf("CorsHeaders() mismatch (-want +got):\n%s", diff)
		}
	})
}

func Test_SetCorsHeaders(t *testing.T) {
	t.Run("adds headers to the response", func(t *testing.T) {
		headers := httpmiddleware.CorsHeaders()

		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		httpmiddleware.SetCorsHeaders(handler).ServeHTTP(w, r)

		if diff := cmp.Diff(headers, w.Result().Header); diff != "" {
			t.Errorf("headers mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("ok on OPTIONS request", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodOptions, "/", nil)

		httpmiddleware.SetCorsHeaders(handler).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})
}
