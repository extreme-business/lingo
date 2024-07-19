package httpmiddleware

import (
	"context"
	"net/http"
)

type Validator interface {
	Validate(ctx context.Context, value string) error
}

func AuthCookie(cookieName string, a Validator, failureRedirect string, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("could not read cookie"))
		}

		if err := a.Validate(r.Context(), cookie.Value); err != nil {
			http.Redirect(w, r, failureRedirect, http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
