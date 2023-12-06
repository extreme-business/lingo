package grpcserver

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func WithCredentials(creds credentials.TransportCredentials) Option {
	return func(c config) config {
		c.ServerOptions = append(c.ServerOptions, grpc.Creds(creds))
		return c
	}
}

func WithReflection() Option {
	return func(c config) config {
		c.Reflection = true
		return c
	}
}
