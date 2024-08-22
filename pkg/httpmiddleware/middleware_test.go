package httpmiddleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/extreme-business/lingo/pkg/httpmiddleware"
)

func TestChain(t *testing.T) {
	// Middleware that adds a header
	addHeaderMiddleware := func(headerName, headerValue string) httpmiddleware.Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add(headerName, headerValue)
				next.ServeHTTP(w, r)
			})
		}
	}

	// Test handler that writes "OK"
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte("OK")); err != nil {
			t.Logf("error writing response: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	tests := []struct {
		name            string
		middlewares     []httpmiddleware.Middleware
		expectedBody    string
		expectedHeaders map[string]string
	}{
		{
			name:            "No middlewares",
			middlewares:     []httpmiddleware.Middleware{},
			expectedBody:    "OK",
			expectedHeaders: map[string]string{},
		},
		{
			name: "One middleware",
			middlewares: []httpmiddleware.Middleware{
				addHeaderMiddleware("X-Test-Header", "test-value"),
			},
			expectedBody: "OK",
			expectedHeaders: map[string]string{
				"X-Test-Header": "test-value",
			},
		},
		{
			name: "Multiple middlewares",
			middlewares: []httpmiddleware.Middleware{
				addHeaderMiddleware("X-Test-Header-1", "value1"),
				addHeaderMiddleware("X-Test-Header-2", "value2"),
			},
			expectedBody: "OK",
			expectedHeaders: map[string]string{
				"X-Test-Header-1": "value1",
				"X-Test-Header-2": "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chainedHandler := httpmiddleware.Chain(tt.middlewares...)(testHandler)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()

			chainedHandler.ServeHTTP(rr, req)

			// Check the response body
			if body := rr.Body.String(); body != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, body)
			}

			// Check the response headers
			for header, expectedValue := range tt.expectedHeaders {
				if value := rr.Header().Get(header); value != expectedValue {
					t.Errorf("expected header %s to be %q, got %q", header, expectedValue, value)
				}
			}
		})
	}
}
