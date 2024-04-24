package httpserver

import (
	"net/http"
	"time"
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
	Headers         http.Header
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

// WithReadTimeout sets the read timeout for the server.
func WithReadTimeout(readTimeout time.Duration) Option {
	return optionFunc(func(c *Config) {
		c.ReadTimeout = readTimeout
	})
}

// WithWriteTimeout sets the write timeout for the server.
func WithWriteTimeout(writeTimeout time.Duration) Option {
	return optionFunc(func(c *Config) {
		c.WriteTimeout = writeTimeout
	})
}

// WithIdleTimeout sets the idle timeout for the server.
func WithIdleTimeout(idleTimeout time.Duration) Option {
	return optionFunc(func(c *Config) {
		c.IdleTimeout = idleTimeout
	})
}

// WithShutdownTimeout sets the shutdown timeout for the server.
func WithShutdownTimeout(shutdownTimeout time.Duration) Option {
	return optionFunc(func(c *Config) {
		c.ShutdownTimeout = shutdownTimeout
	})
}

// WithTLS sets the certificate and key files for the server.
func WithTLS(certFile string, keyFile string) Option {
	return optionFunc(func(c *Config) {
		c.CertFile = certFile
		c.KeyFile = keyFile
	})
}

// WithHeaders sets the headers for the server that will be used in the responses.
func WithHeaders(headers http.Header) Option {
	return optionFunc(func(c *Config) {
		c.Headers = headers
	})
}
