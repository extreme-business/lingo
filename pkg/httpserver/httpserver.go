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

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			fmt.Printf("failed to shutdown http server: %v\n", err)
		}
	}()

	if err := s.httpServer.ListenAndServeTLS(s.certFile, s.keyFile); err != nil {
		if err == http.ErrServerClosed || errors.Is(err, context.Canceled) {
			return nil
		}

		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
