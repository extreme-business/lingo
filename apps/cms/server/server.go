package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/extreme-business/lingo/apps/account/domain"
	"github.com/extreme-business/lingo/apps/cms/account"
	"github.com/extreme-business/lingo/apps/cms/cookie"
	"github.com/extreme-business/lingo/apps/cms/views"
	"github.com/extreme-business/lingo/pkg/httpmiddleware"
	"github.com/extreme-business/lingo/pkg/httpserver"
	"github.com/google/uuid"
)

const (
	readTimeout     = 5 * time.Second
	writeTimeout    = 10 * time.Second
	idleTimeout     = 15 * time.Second
	shutdownTimeout = 5 * time.Second
)

const (
	AccessTokenDuration  = 24 * time.Hour
	RefreshTokenDuration = 7 * 24 * time.Hour
)

type Registration struct {
	Email    string
	Password string
}

// AccountManager is the interface for authenticating users.
type AccountManager interface {
	Authenticate(ctx context.Context, email, password string) (*account.SuccessResponse, error)
	Register(ctx context.Context, r account.Registration) error
}

// New creates a new Server instance.
func New(
	addr string,
	accountManager AccountManager,
	authMiddleware httpmiddleware.Middleware,
) *httpserver.Server {
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/", homeHandler)

	mux := http.NewServeMux()
	mux.Handle("/", authMiddleware(adminMux))

	mux.HandleFunc("/login", loginHandler(time.Now, accountManager))
	mux.HandleFunc("/register", registerHandler(accountManager))

	return httpserver.New(
		httpserver.WithAddr(addr),
		httpserver.WithHandler(mux),
		httpserver.WithTimeouts(httpserver.Timeouts{
			ReadTimeout:     readTimeout,
			WriteTimeout:    writeTimeout,
			IdleTimeout:     idleTimeout,
			ShutdownTimeout: shutdownTimeout,
		}),
		httpserver.WithMiddleware(httpmiddleware.SetCorsHeaders),
	)
}

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := views.UserList(w, []*domain.User{
		{
			ID: uuid.Max,
		},
	})
	if err != nil {
		if err = views.Error(w, err.Error()); err != nil {
			slog.Error(err.Error())
			http.Error(w, "could not render view", http.StatusInternalServerError)
		}
	}
}

const loginTemplate = `
<!DOCTYPE html>
<html>
<head>
	<title>Home</title>
</head>
<body>
	<h1>Welcome to the Home Page</h1>
	<div>
		<form action="/login" method="post">
			<label for="email">email:</label>
			<input type="email" id="email" name="email">
			<label for="password">Password:</label>
			<input type="password" id="password" name="password">
			<button type="submit">Login</button>
		</form>
	</div>
</body>
</html>
`

func loginHandler(c func() time.Time, a AccountManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		if r.Method == http.MethodPost {
			r.ParseForm()
			email := r.Form.Get("email")
			password := r.Form.Get("password")

			s, err := a.Authenticate(r.Context(), email, password)
			if err != nil {
				if err = views.Error(w, err.Error()); err != nil {
					slog.Error(err.Error())
					http.Error(w, "could not render view", http.StatusInternalServerError)
				}
				return
			}

			now := c()
			cookie.SetAccessToken(w, s.AccessToken, now.Add(AccessTokenDuration))
			cookie.SetRefreshToken(w, s.RefreshToken, now.Add(RefreshTokenDuration))

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		fmt.Fprint(w, loginTemplate)
	}
}
