package httpserver_test

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/extreme-business/lingo/pkg/httpserver"
)

func TestWithAddr(t *testing.T) {
	t.Run("WithAddr", func(t *testing.T) {
		addr := "localhost:8080"
		opt := httpserver.WithAddr(addr)

		c := &httpserver.Config{}
		c.Apply(opt)

		if c.Addr != addr {
			t.Errorf("expected %s, got %s", addr, c.Addr)
		}
	})
}

func TestWithHandler(t *testing.T) {
	t.Run("WithHandler", func(t *testing.T) {
		handler := http.DefaultServeMux
		opt := httpserver.WithHandler(handler)

		c := &httpserver.Config{}
		c.Apply(opt)

		if c.Handler != handler {
			t.Errorf("expected %v, got %v", handler, c.Handler)
		}
	})
}

func TestWithReadTimeout(t *testing.T) {
	t.Run("WithReadTimeout", func(t *testing.T) {
		readTimeout := 5 * time.Second
		opt := httpserver.WithReadTimeout(readTimeout)

		c := &httpserver.Config{}
		c.Apply(opt)

		if c.ReadTimeout != readTimeout {
			t.Errorf("expected %v, got %v", readTimeout, c.ReadTimeout)
		}
	})
}

func TestWithWriteTimeout(t *testing.T) {
	t.Run("WithWriteTimeout", func(t *testing.T) {
		writeTimeout := 10 * time.Second
		opt := httpserver.WithWriteTimeout(writeTimeout)

		c := &httpserver.Config{}
		c.Apply(opt)

		if c.WriteTimeout != writeTimeout {
			t.Errorf("expected %v, got %v", writeTimeout, c.WriteTimeout)
		}
	})
}

func TestWithIdleTimeout(t *testing.T) {
	t.Run("WithIdleTimeout", func(t *testing.T) {
		idleTimeout := 15 * time.Second
		opt := httpserver.WithIdleTimeout(idleTimeout)

		c := &httpserver.Config{}
		c.Apply(opt)

		if c.IdleTimeout != idleTimeout {
			t.Errorf("expected %v, got %v", idleTimeout, c.IdleTimeout)
		}
	})
}

func TestWithShutdownTimeout(t *testing.T) {
	t.Run("WithShutdownTimeout", func(t *testing.T) {
		shutdownTimeout := 5 * time.Second
		opt := httpserver.WithShutdownTimeout(shutdownTimeout)

		c := &httpserver.Config{}
		c.Apply(opt)

		if c.ShutdownTimeout != shutdownTimeout {
			t.Errorf("expected %v, got %v", shutdownTimeout, c.ShutdownTimeout)
		}
	})
}

func TestWithTLS(t *testing.T) {
	t.Run("WithTLS", func(t *testing.T) {
		certFile := "certFile"
		keyFile := "keyFile"
		opt := httpserver.WithTLS(certFile, keyFile)

		c := &httpserver.Config{}
		c.Apply(opt)

		if c.CertFile != certFile {
			t.Errorf("expected %s, got %s", certFile, c.CertFile)
		}

		if c.KeyFile != keyFile {
			t.Errorf("expected %s, got %s", keyFile, c.KeyFile)
		}
	})
}

func TestWithHeaders(t *testing.T) {
	t.Run("WithHeaders", func(t *testing.T) {
		headers := httpserver.CorsHeaders()
		opt := httpserver.WithHeaders(headers)

		c := &httpserver.Config{}
		c.Apply(opt)

		if !reflect.DeepEqual(c.Headers, headers) {
			t.Errorf("expected %v, got %v", headers, c.Headers)
		}
	})
}
