package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/extreme-business/lingo/apps/cms/account"
)

func registerHandler(a AccountManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				slog.ErrorContext(r.Context(), "failed to parse form", slog.String("error", err.Error()))
				http.Error(w, "failed to parse form", http.StatusBadRequest)
				return
			}

			email := r.Form.Get("email")
			password := r.Form.Get("password")

			err := a.Register(r.Context(), account.Registration{
				Email:    email,
				Password: password,
			})

			if err != nil {
				slog.ErrorContext(r.Context(), "failed to create user", slog.String("error", err.Error()))
			}

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
        <!DOCTYPE html>
        <html>
        <head>
            <title>Home</title>
        </head>
        <body>
            <h1>Welcome to the Home Page</h1>
            <nav>
                <a href="/">Home</a> |
                <a href="/about">About</a> |
                <a href="/contact">Contact</a>
            </nav>
            <div>
                <form action="/login" method="post">
                    <label for="organization">Choose a flavor:</label>
                    <input list="organizations" id="organizations name="organization" />
                    <datalist id="organizations">
                        <option value="Chocolate"></option>
                        <option value="Coconut"></option>
                        <option value="Mint"></option>
                        <option value="Strawberry"></option>
                        <option value="Vanilla"></option>
                    </datalist>
                    <label for="username">Username:</label>
                    <input type="text" id="username" name="username">
                    <label for="password">Password:</label>
                    <input type="password" id="password" name="password">
                    <button type="submit">Register</button>
                </form>
            </div>
        </body>
        </html>
    `)
	}
}
