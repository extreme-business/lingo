package httpmiddleware

import (
	"context"
	"log/slog"
	"net/http"
	"slices"
)

type TokenValidator interface {
	Validate(ctx context.Context, token string) error
}

func AuthCookie(cookieName string, a TokenValidator, failureRedirect string, excludePathsAndMethods map[string][]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if methods, ok := excludePathsAndMethods[r.URL.Path]; ok {
				if slices.Contains(methods, r.Method) {
					next.ServeHTTP(w, r)
					return
				}
			}

			cookie, err := r.Cookie(cookieName)
			if err != nil {
				slog.InfoContext(r.Context(), "cookie not found")
				http.Redirect(w, r, failureRedirect, http.StatusSeeOther)
				return
			}

			if err := a.Validate(r.Context(), cookie.Value); err != nil {
				slog.InfoContext(r.Context(), "cookie validation failed", slog.String("error", err.Error()))
				http.Redirect(w, r, failureRedirect, http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
