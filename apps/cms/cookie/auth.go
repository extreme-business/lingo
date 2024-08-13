package cookie

import (
	"net/http"
	"time"
)

func AccessToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func SetAccessToken(w http.ResponseWriter, token string, expires time.Time) {
	Set(w, "access_token", token, expires, true, "/")
}

func RefreshToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func SetRefreshToken(w http.ResponseWriter, token string, expires time.Time) {
	Set(w, "refresh_token", token, expires, true, "/")
}
