package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/extreme-business/lingo/apps/account/domain"
	"github.com/extreme-business/lingo/apps/cms/app"
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

// New creates a new Server instance.
func New(
	logger *slog.Logger,
	addr string,
	app *app.App,
	authMiddleware httpmiddleware.Middleware,
) (*httpserver.Server, error) {
	vw, err := views.New()
	if err != nil {
		logger.Error(err.Error())
		return nil, fmt.Errorf("failed to create views writer: %w", err)
	}

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/", homeHandler(vw))

	mux := http.NewServeMux()
	mux.Handle("/", authMiddleware(adminMux))

	mux.HandleFunc("/login", loginHandler(vw, logger, time.Now, app))
	mux.HandleFunc("POST /logout", logoutHandler(vw))
	mux.HandleFunc("/register", registerHandler(vw, logger, app))

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
	), nil
}

func homeHandler(vw *views.Writer) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		err := vw.UserList(w, []*domain.User{
			{
				ID: uuid.Max,
			},
		})
		if err != nil {
			if err = vw.Error(w, err.Error()); err != nil {
				slog.Error(err.Error())
				http.Error(w, "could not render view", http.StatusInternalServerError)
			}
		}
	}
}

func loginHandler(
	vw *views.Writer,
	logger *slog.Logger,
	c func() time.Time,
	app *app.App,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				logger.Error(err.Error())
				http.Error(w, "failed to parse form", http.StatusBadRequest)
				return
			}

			email := r.Form.Get("email")
			password := r.Form.Get("password")

			s, err := app.AuthenticateUser(r.Context(), email, password)
			if err != nil {
				logger.Error(err.Error())
				if err = vw.Error(w, "could not login"); err != nil {
					logger.Error(err.Error())
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

		if err := vw.Login(w); err != nil {
			logger.Error(err.Error())
			http.Error(w, "could not render view", http.StatusInternalServerError)
		}
	}
}

func logoutHandler(vw *views.Writer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		vw.Logout(w)
	}
}
