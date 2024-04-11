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

type Config struct {
	Addr            string
	Handler         http.Handler
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	CertFile        string // CertFile is the path to the certificate file
	KeyFile         string // KeyFile is the path to the key file
	Cors            bool
}

type Server struct {
	httpServer        *http.Server
	shutdownTimeout   time.Duration
	certFile, keyFile string
}

func New(c Config) *Server {
	httpServer := &http.Server{
		Addr:         c.Addr,
		Handler:      c.Handler,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		IdleTimeout:  c.IdleTimeout,
	}

	if c.Cors {
		httpServer.Handler = CorsMiddleware(httpServer.Handler)
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
		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		if err := s.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
		return ctx.Err()
	case err := <-errCh:
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return fmt.Errorf("failed to serve: %w", err)
	}
}
