package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	ErrServerAlreadyRunning = errors.New("server is already running")
	ErrListenerNotSet       = errors.New("listener is not set")
)

type Config struct {
	Listener        net.Listener                  // Listener
	ServerRegisters []func(*grpc.Server)          // Server registers
	ServerOptions   []grpc.ServerOption           // Server options
	Reflection      bool                          // Enable reflection
	reflectionFunc  func(s reflection.GRPCServer) // Reflection function, mostly for testing
}

// grpcServer is an interface that wraps the Serve and GracefulStop methods.
type grpcServer interface {
	Serve(net.Listener) error
	GracefulStop()
}

type Server struct {
	lis     net.Listener
	Serv    grpcServer
	running bool
}

func New(c Config) *Server {
	grpcServer := grpc.NewServer(c.ServerOptions...)

	if c.Reflection {
		if c.reflectionFunc == nil {
			c.reflectionFunc = reflection.Register
		}

		c.reflectionFunc(grpcServer)
	}

	for _, register := range c.ServerRegisters {
		register(grpcServer)
	}

	return &Server{
		lis:  c.Listener,
		Serv: grpcServer,
	}
}

// Running returns true if the server is running.
func (s *Server) Running() bool { return s.running }

// Serve starts the grpc server.
func (s *Server) Serve(ctx context.Context) error {
	if s.running {
		return ErrServerAlreadyRunning
	}

	s.running = true
	defer func() {
		s.running = false
	}()

	if s.lis == nil {
		return ErrListenerNotSet
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- s.Serv.Serve(s.lis)
	}()

	select {
	case <-ctx.Done():
		s.Serv.GracefulStop()
		return ctx.Err()
	case err := <-errChan:
		if err == nil || errors.Is(err, grpc.ErrServerStopped) {
			return nil
		}

		return fmt.Errorf("failed to serve: %w", err)
	}
}
