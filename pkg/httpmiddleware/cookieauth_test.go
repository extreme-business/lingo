package httpmiddleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/extreme-business/lingo/pkg/httpmiddleware"
)

// MockValidator is a mock implementation of the Validator interface for testing.
type MockValidator struct {
	ValidateFunc func(ctx context.Context, value string) error
}

func (m *MockValidator) Validate(ctx context.Context, value string) error {
	return m.ValidateFunc(ctx, value)
}

func TestAuthCookie(t *testing.T) {
	tests := []struct {
		name           string
		cookieName     string
		cookieValue    string
		validateFunc   func(ctx context.Context, value string) error
		expectedStatus int
		expectedURL    string
	}{
		{
			name:        "Valid cookie",
			cookieName:  "auth",
			cookieValue: "valid_token",
			validateFunc: func(ctx context.Context, value string) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedURL:    "",
		},
		{
			name:        "Missing cookie",
			cookieName:  "auth",
			cookieValue: "",
			validateFunc: func(ctx context.Context, value string) error {
				return nil
			},
			expectedStatus: http.StatusSeeOther,
			expectedURL:    "/failure",
		},
		{
			name:        "Invalid cookie",
			cookieName:  "auth",
			cookieValue: "invalid_token",
			validateFunc: func(ctx context.Context, value string) error {
				return errors.New("invalid token")
			},
			expectedStatus: http.StatusSeeOther,
			expectedURL:    "/failure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &MockValidator{
				ValidateFunc: tt.validateFunc,
			}

			handler := httpmiddleware.AuthCookie(tt.cookieName, validator, "/failure")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.cookieValue != "" {
				req.AddCookie(&http.Cookie{Name: tt.cookieName, Value: tt.cookieValue})
			}

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if status := rr.Result().StatusCode; status != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, status)
			}

			if tt.expectedURL != "" {
				location, err := rr.Result().Location()
				if err != nil {
					t.Fatalf("expected redirect, got error: %v", err)
				}
				if location.String() != tt.expectedURL {
					t.Errorf("expected redirect to %v, got %v", tt.expectedURL, location.String())
				}
			}
		})
	}
}
