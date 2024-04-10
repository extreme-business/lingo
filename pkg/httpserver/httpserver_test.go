package httpserver

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		c Config
	}
	tests := []struct {
		name string
		args args
		want *Server
	}{
		{
			name: "TestNew",
			args: args{
				c: Config{
					Addr:            "localhost:8080",
					Handler:         http.DefaultServeMux,
					ReadTimeout:     5 * time.Second,
					WriteTimeout:    10 * time.Second,
					IdleTimeout:     15 * time.Second,
					ShutdownTimeout: 5 * time.Second,
					CertFile:        "cert.pem",
					KeyFile:         "key.pem",
				},
			},
			want: &Server{
				httpServer: &http.Server{
					Addr:         "localhost:8080",
					Handler:      http.DefaultServeMux,
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  15 * time.Second,
				},
				shutdownTimeout: 5 * time.Second,
				certFile:        "cert.pem",
				keyFile:         "key.pem",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.c)
			if got.httpServer.ReadTimeout != tt.want.httpServer.ReadTimeout {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
			if got.httpServer.WriteTimeout != tt.want.httpServer.WriteTimeout {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
			if got.httpServer.IdleTimeout != tt.want.httpServer.IdleTimeout {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
			if got.shutdownTimeout != tt.want.shutdownTimeout {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
			if got.certFile != tt.want.certFile {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
			if got.keyFile != tt.want.keyFile {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
			if got.httpServer.Addr != tt.want.httpServer.Addr {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
			if got.httpServer.Handler != tt.want.httpServer.Handler {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
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
			err: ErrCertFilesNotSet,
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
			err: ErrCertFilesNotSet,
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
			err: ErrCertFilesNotSet,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				httpServer:      tt.fields.httpServer,
				shutdownTimeout: tt.fields.shutdownTimeout,
				certFile:        tt.fields.certFile,
				keyFile:         tt.fields.keyFile,
			}

			err := s.Serve(tt.args.ctx)
			if !errors.Is(tt.err, err) {
				t.Errorf("Serve() = %v, want %v", err, tt.err)
			}
		})
	}
}
