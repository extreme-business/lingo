package cookie_test

import (
	"net/http"
	"testing"

	"github.com/extreme-business/lingo/apps/cms/cookie"
)

func TestGetAccessToken(t *testing.T) {
	t.Run("should return access token from request", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "access",
		})

		token, err := cookie.AccessToken(r)
		if err != nil {
			t.Error("expected no error")
		}

		if token != "access" {
			t.Errorf("expected access token to be 'access', got %s", token)
		}
	})

	t.Run("should return error if access token is not found", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		if _, err := cookie.AccessToken(r); err == nil {
			t.Error("expected error")
		}
	})
}

func TestGetRefreshToken(t *testing.T) {
	t.Run("should return refresh token from request", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.AddCookie(&http.Cookie{
			Name:  "refresh_token",
			Value: "refresh",
		})

		token, err := cookie.RefreshToken(r)
		if err != nil {
			t.Error("expected no error")
		}

		if token != "refresh" {
			t.Errorf("expected refresh token to be 'refresh', got %s", token)
		}
	})

	t.Run("should return error if refresh token is not found", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		if _, err := cookie.RefreshToken(r); err == nil {
			t.Error("expected error")
		}
	})
}
