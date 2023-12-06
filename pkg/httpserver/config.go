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

func WithAddr(addr string) Option {
	return optionFunc(func(c *Config) {
		c.Addr = addr
	})
}

func WithHandler(handler http.Handler) Option {
	return optionFunc(func(c *Config) {
		c.Handler = handler
	})
}

func WithReadTimeout(readTimeout time.Duration) Option {
	return optionFunc(func(c *Config) {
		c.ReadTimeout = readTimeout
	})
}

func WithWriteTimeout(writeTimeout time.Duration) Option {
	return optionFunc(func(c *Config) {
		c.WriteTimeout = writeTimeout
	})
}

func WithIdleTimeout(idleTimeout time.Duration) Option {
	return optionFunc(func(c *Config) {
		c.IdleTimeout = idleTimeout
	})
}

func WithShutdownTimeout(shutdownTimeout time.Duration) Option {
	return optionFunc(func(c *Config) {
		c.ShutdownTimeout = shutdownTimeout
	})
}

func WithTLS(certFile string, keyFile string) Option {
	return optionFunc(func(c *Config) {
		c.CertFile = certFile
		c.KeyFile = keyFile
	})
}

func WithHeaders(headers http.Header) Option {
	return optionFunc(func(c *Config) {
		c.Headers = headers
	})
}
