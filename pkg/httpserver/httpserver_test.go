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
				httpserver.WithTimeouts(httpserver.Timeouts{
					ReadTimeout:     5 * time.Second,
					WriteTimeout:    10 * time.Second,
					IdleTimeout:     15 * time.Second,
					ShutdownTimeout: 5 * time.Second,
				}),
			)

			if s == nil {
				t.Error("New() = nil, want not nil")
			}
		})
	})
}

func TestServer_ServeTLS(t *testing.T) {
	type fields struct {
		httpServer      *http.Server
		shutdownTimeout time.Duration
	}
	type args struct {
		ctx      context.Context
		certFile string
		keyFile  string
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
			},
			args: args{
				ctx:      context.Background(),
				certFile: "",
				keyFile:  "",
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
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
				certFile: "./testdata/test.crt",
				keyFile:  "./testdata/test.key",
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
			},
			args: args{
				ctx:      context.Background(),
				certFile: "",
				keyFile:  "./testdata/test.key",
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
			},
			args: args{
				ctx:      context.Background(),
				certFile: "./testdata/test.crt",
				keyFile:  "",
			},
			err: httpserver.ErrCertFilesNotSet,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := httpserver.New(
				httpserver.WithAddr(tt.fields.httpServer.Addr),
				httpserver.WithHandler(tt.fields.httpServer.Handler),
				httpserver.WithTimeouts(httpserver.Timeouts{
					ReadTimeout:     tt.fields.httpServer.ReadTimeout,
					WriteTimeout:    tt.fields.httpServer.WriteTimeout,
					IdleTimeout:     tt.fields.httpServer.IdleTimeout,
					ShutdownTimeout: tt.fields.shutdownTimeout,
				}),
			)

			err := s.ServeTLS(tt.args.ctx, tt.args.certFile, tt.args.keyFile)
			if !errors.Is(tt.err, err) {
				t.Errorf("Serve() = %v, want %v", err, tt.err)
			}
		})
	}
}
