package httpserver_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/extreme-business/lingo/pkg/httpserver"
)

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		t.Run("should return a new server", func(t *testing.T) {
			s := httpserver.New(
				httpserver.WithAddr("localhost:8080"),
				httpserver.WithHandler(http.DefaultServeMux),
				httpserver.WithReadTimeout(5*time.Second),
				httpserver.WithWriteTimeout(10*time.Second),
				httpserver.WithIdleTimeout(15*time.Second),
				httpserver.WithShutdownTimeout(5*time.Second),
				httpserver.WithTLS("./testdata/test.crt", "./testdata/test.key"),
			)

			if s == nil {
				t.Error("New() = nil, want not nil")
			}
		})
	})
}

func TestServer_Serve(t *testing.T) {
	type fields struct {
		httpServer      *http.Server
		shutdownTimeout time.Duration
		certFile        string
		keyFile         string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		err    error
	}{
		{
			name: "should return error when certFile and keyFile are not set",
			fields: fields{
				httpServer: &http.Server{
					Addr:         "localhost:8081",
					Handler:      http.DefaultServeMux,
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  15 * time.Second,
				},
				shutdownTimeout: 5 * time.Second,
				certFile:        "",
				keyFile:         "",
			},
			args: args{
				ctx: context.Background(),
			},
			err: httpserver.ErrCertFilesNotSet,
		},
		{
			name: "should not return error when certFile and keyFile are set",
			fields: fields{
				httpServer: &http.Server{
					Addr:         "localhost:8082",
					Handler:      http.DefaultServeMux,
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  15 * time.Second,
				},
				shutdownTimeout: 5 * time.Second,
				certFile:        "./testdata/test.crt",
				keyFile:         "./testdata/test.key",
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
			},
			err: nil,
		},
		{
			name: "should return error when certFile is not set",
			fields: fields{
				httpServer: &http.Server{
					Addr:         "localhost:8083",
					Handler:      http.DefaultServeMux,
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  15 * time.Second,
				},
				shutdownTimeout: 5 * time.Second,
				certFile:        "",
				keyFile:         "./testdata/test.key",
			},
			args: args{
				ctx: context.Background(),
			},
			err: httpserver.ErrCertFilesNotSet,
		},
		{
			name: "should return error when keyFile is not set",
			fields: fields{
				httpServer: &http.Server{
					Addr:         "localhost:8084",
					Handler:      http.DefaultServeMux,
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  15 * time.Second,
				},
				shutdownTimeout: 5 * time.Second,
				certFile:        "./testdata/test.crt",
				keyFile:         "",
			},
			args: args{
				ctx: context.Background(),
			},
			err: httpserver.ErrCertFilesNotSet,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := httpserver.New(
				httpserver.WithAddr(tt.fields.httpServer.Addr),
				httpserver.WithHandler(tt.fields.httpServer.Handler),
				httpserver.WithReadTimeout(tt.fields.httpServer.ReadTimeout),
				httpserver.WithWriteTimeout(tt.fields.httpServer.WriteTimeout),
				httpserver.WithIdleTimeout(tt.fields.httpServer.IdleTimeout),
				httpserver.WithShutdownTimeout(tt.fields.shutdownTimeout),
				httpserver.WithTLS(tt.fields.certFile, tt.fields.keyFile),
			)

			err := s.Serve(tt.args.ctx)
			if !errors.Is(tt.err, err) {
				t.Errorf("Serve() = %v, want %v", err, tt.err)
			}
		})
	}
}
