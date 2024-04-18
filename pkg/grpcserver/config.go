package grpcserver

import (
	"net"

	"google.golang.org/grpc"
)

type Config struct {
	Lis             net.Listener                  // Listener
	Address         string                        // Address
	ServerRegisters []func(grpc.ServiceRegistrar) // Server registers
	Reflection      bool                          // Enable reflection
	GrpcServer      grpcServer                    // Grpc server
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

// WithListener sets the listener for the server.
func WithListener(l net.Listener) Option {
	return optionFunc(func(c *Config) {
		c.Lis = l
	})
}

func WithServiceRegistrars(r []func(grpc.ServiceRegistrar)) Option {
	return optionFunc(func(c *Config) {
		c.ServerRegisters = r
	})
}

func WithReflection() Option {
	return optionFunc(func(c *Config) {
		c.Reflection = true
	})
}

// WithAddress sets the address of the server and is used if no listener is set.
func WithAddress(a string) Option {
	return optionFunc(func(c *Config) {
		c.Address = a
	})
}

// WithGrpcServer sets the grpc server. For testing purposes.
func WithGrpcServer(s grpcServer) Option {
	return optionFunc(func(c *Config) {
		c.GrpcServer = s
	})
}
