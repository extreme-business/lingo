package httpmiddleware

import (
	"net/http"
	"strings"
)

// CorsHeaders returns a copy of the default CORS headers.
func CorsHeaders() http.Header {
	return http.Header{
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
}

// SetHeaders is a middleware that adds headers to the response.
func SetCorsHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, values := range CorsHeaders() {
			for _, value := range values {
				w.Header().Set(key, value)
			}
		}

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
