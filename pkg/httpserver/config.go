package httpserver

import (
	"net/http"
	"time"

	"github.com/extreme-business/lingo/pkg/httpmiddleware"
)

type Config struct {
	Addr       string
	Middleware []httpmiddleware.Middleware
	Handler    http.Handler
	// timeouts
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration

	CertFile string // CertFile is the path to the certificate file
	KeyFile  string // KeyFile is the path to the key file
}

func (c *Config) Apply(opts ...Option) {
	for _, o := range opts {
		o.apply(c)
	}
}

type Option interface {
	apply(*Config)
}

type optionFunc func(*Config)

func (f optionFunc) apply(c *Config) {
	f(c)
}

// WithAddr sets the address the server will listen on.
func WithAddr(addr string) Option {
	return optionFunc(func(c *Config) {
		c.Addr = addr
	})
}

// WithHandler sets the handler the server will use.
func WithHandler(handler http.Handler) Option {
	return optionFunc(func(c *Config) {
		c.Handler = handler
	})
}

type Timeouts struct {
	ReadTimeout     time.Duration // ReadTimeout is the maximum duration for reading the entire request, including the body.
	WriteTimeout    time.Duration // WriteTimeout is the maximum duration before timing out writes of the response.
	IdleTimeout     time.Duration // IdleTimeout is the maximum amount of time to wait for the next request when keep-alives are enabled.
	ShutdownTimeout time.Duration // ShutdownTimeout is the maximum duration before shutting down the server.
}

func WithTimeouts(t Timeouts) Option {
	return optionFunc(func(c *Config) {
		c.ReadTimeout = t.ReadTimeout
		c.WriteTimeout = t.WriteTimeout
		c.IdleTimeout = t.IdleTimeout
		c.ShutdownTimeout = t.ShutdownTimeout
	})
}

// WithHeaders sets the headers for the server that will be used in the responses.
func WithMiddleware(m ...httpmiddleware.Middleware) Option {
	return optionFunc(func(c *Config) {
		c.Middleware = append(c.Middleware, m...)
	})
}
