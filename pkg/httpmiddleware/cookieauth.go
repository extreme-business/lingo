package httpmiddleware

import (
	"context"
	"net/http"
)

type TokenValidator interface {
	Validate(ctx context.Context, token string) error
}

func AuthCookie(cookieName string, a TokenValidator, failureRedirect string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(cookieName)
			if err != nil {
				http.Redirect(w, r, failureRedirect, http.StatusSeeOther)
				return
			}

			if err := a.Validate(r.Context(), cookie.Value); err != nil {
				http.Redirect(w, r, failureRedirect, http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
