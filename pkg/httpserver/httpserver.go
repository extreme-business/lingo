package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrCertFilesNotSet = errors.New("certFile and keyFile must be set")
)

type Server struct {
	httpServer        *http.Server
	shutdownTimeout   time.Duration
	certFile, keyFile string
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

	if c.Headers != nil && len(c.Headers) > 0 {
		httpServer.Handler = HeadersMiddleware(httpServer.Handler, c.Headers)
	}

	return &Server{
		httpServer:      httpServer,
		shutdownTimeout: c.ShutdownTimeout,
		certFile:        c.CertFile,
		keyFile:         c.KeyFile,
	}
}

func (s *Server) Serve(ctx context.Context) error {
	if s.certFile == "" || s.keyFile == "" {
		return ErrCertFilesNotSet
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.httpServer.ListenAndServeTLS(s.certFile, s.keyFile)
	}()

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
