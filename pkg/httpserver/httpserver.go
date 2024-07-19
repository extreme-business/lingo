package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/extreme-business/lingo/pkg/httpmiddleware"
)

var (
	ErrCertFilesNotSet      = errors.New("certFile and keyFile must be set")
	ErrServerAlreadyStarted = errors.New("server already started")
)

type Server struct {
	httpServer      *http.Server
	shutdownTimeout time.Duration
	started         bool
}

func New(options ...Option) *Server {
	c := &Config{}
	c.Apply(options...)

	httpServer := &http.Server{
		Addr:         c.Addr,
		Handler:      c.Handler,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		IdleTimeout:  c.IdleTimeout,
	}

	if c.Middleware != nil {
		httpServer.Handler = httpmiddleware.Chain(c.Middleware...)(httpServer.Handler)
	}

	return &Server{
		httpServer:      httpServer,
		shutdownTimeout: c.ShutdownTimeout,
	}
}

func (s *Server) ServeTLS(ctx context.Context, certFile, keyFile string) error {
	if s.started {
		return ErrServerAlreadyStarted
	}

	if certFile == "" || keyFile == "" {
		return ErrCertFilesNotSet
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.httpServer.ListenAndServeTLS(certFile, keyFile)
	}()

	return s.handleErr(ctx, errCh)
}

func (s *Server) Serve(ctx context.Context) error {
	if s.started {
		return ErrServerAlreadyStarted
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.httpServer.ListenAndServe()
	}()

	return s.handleErr(ctx, errCh)
}

func (s *Server) handleErr(ctx context.Context, errCh chan error) error {
	s.started = true

	select {
	case <-ctx.Done():
		ctxt, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		if err := s.httpServer.Shutdown(ctxt); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
		return ctxt.Err()
	case err := <-errCh:
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return fmt.Errorf("failed to serve: %w", err)
	}
}
