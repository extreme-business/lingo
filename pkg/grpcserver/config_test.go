package grpcserver_test

import (
	"net"
	"testing"

	"github.com/extreme-business/lingo/pkg/grpcserver"
	"google.golang.org/grpc"
)

func TestWithListener(t *testing.T) {
	t.Run("WithListener", func(t *testing.T) {
		lis := &net.TCPListener{}
		opt := grpcserver.WithListener(lis)

		c := &grpcserver.Config{}
		c.Apply(opt)

		if c.Lis != lis {
			t.Errorf("expected %v, got %v", lis, c.Lis)
		}
	})
}

func TestWithServiceRegistrars(t *testing.T) {
	t.Run("registrar should be registered", func(t *testing.T) {
		r := func(grpc.ServiceRegistrar) {}
		opt := grpcserver.WithServiceRegistrar(r)

		c := &grpcserver.Config{}
		c.Apply(opt)

		if c.ServiceRegistrar == nil {
			t.Errorf("expected not nil, got nil")
		}
	})

	t.Run("registrar should not be registered", func(t *testing.T) {
		opt := grpcserver.WithServiceRegistrar(nil)

		c := &grpcserver.Config{}
		c.Apply(opt)

		if c.ServiceRegistrar != nil {
			t.Errorf("expected nil, got not nil")
		}
	})
}

func TestWithReflection(t *testing.T) {
	t.Run("WithReflection", func(t *testing.T) {
		opt := grpcserver.WithReflection()

		c := &grpcserver.Config{}
		c.Apply(opt)

		if !c.Reflection {
			t.Errorf("expected true, got false")
		}
	})
}

func TestWithAddress(t *testing.T) {
	t.Run("WithAddress", func(t *testing.T) {
		addr := "localhost:8080"
		opt := grpcserver.WithAddress(addr)

		c := &grpcserver.Config{}
		c.Apply(opt)

		if c.Address != addr {
			t.Errorf("expected %s, got %s", addr, c.Address)
		}
	})
}

func TestWithGrpcServer(t *testing.T) {
	t.Run("WithGrpcServer", func(t *testing.T) {
		s := &grpc.Server{}
		opt := grpcserver.WithGrpcServer(s)

		c := &grpcserver.Config{}
		c.Apply(opt)

		if c.GrpcServer != s {
			t.Errorf("expected %v, got %v", s, c.GrpcServer)
		}
	})
}
