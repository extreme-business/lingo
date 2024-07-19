package cookie

import (
	"net/http"
	"time"
)

func Set(w http.ResponseWriter, name, value string, expires time.Time, httpOnly bool, path string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expires,
		HttpOnly: httpOnly,
		Path:     path,
	})
}
