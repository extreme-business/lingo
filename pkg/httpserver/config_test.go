package httpserver_test

import (
	"net/http"
	"reflect"
	"testing"

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

func TestWithTimeouts(t *testing.T) {
	t.Run("WithTimeouts", func(t *testing.T) {
		timeouts := httpserver.Timeouts{
			ReadTimeout:     5,
			WriteTimeout:    10,
			IdleTimeout:     15,
			ShutdownTimeout: 5,
		}
		opt := httpserver.WithTimeouts(timeouts)

		c := &httpserver.Config{}
		c.Apply(opt)

		if !reflect.DeepEqual(c.ReadTimeout, timeouts.ReadTimeout) {
			t.Errorf("expected %v, got %v", timeouts.ReadTimeout, c.ReadTimeout)
		}
		if !reflect.DeepEqual(c.WriteTimeout, timeouts.WriteTimeout) {
			t.Errorf("expected %v, got %v", timeouts.WriteTimeout, c.WriteTimeout)
		}
		if !reflect.DeepEqual(c.IdleTimeout, timeouts.IdleTimeout) {
			t.Errorf("expected %v, got %v", timeouts.IdleTimeout, c.IdleTimeout)
		}
		if !reflect.DeepEqual(c.ShutdownTimeout, timeouts.ShutdownTimeout) {
			t.Errorf("expected %v, got %v", timeouts.ShutdownTimeout, c.ShutdownTimeout)
		}
	})
}
