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

// grpcServer is an interface that wraps the Serve and GracefulStop methods.
type grpcServer interface {
	grpc.ServiceRegistrar
	reflection.GRPCServer
	Serve(net.Listener) error
	GracefulStop()
}

type Server struct {
	lis        net.Listener
	grpcServer grpcServer
	running    bool
	address    string
}

// New creates a new grpc server.
//   - If no listener is set, the server will create a new listener on the address.
//   - If no grpc server is set, the server will create a new grpc server.
//   - If reflection is enabled, the server will register the reflection service.
//   - If server registers are set, the server will register the services.
//
// Example:
//
//	grpcserver.New(
//		grpcserver.WithGrpcServer(grpc.NewServer(grpc.Creds(creds))),
//		grpcserver.WithAddress(fmt.Sprintf(":%d", grpcPort)),
//		grpcserver.WithServiceRegistrars(serviceRegistrars),
//	)
//
// Example with listener:
//
//	grpcserver.New(
//		grpcserver.WithGrpcServer(grpc.NewServer(grpc.Creds(creds))),
//		grpcserver.WithLIstener(lis),
//		grpcserver.WithServiceRegistrars(serviceRegistrars),
//	)
func New(options ...Option) *Server {
	c := &Config{}
	c.Apply(options...)

	if c.GrpcServer == nil {
		c.GrpcServer = grpc.NewServer()
	}

	if c.Reflection {
		reflection.Register(c.GrpcServer)
	}

	for _, register := range c.ServerRegisters {
		register(c.GrpcServer)
	}

	var address string
	if c.Lis != nil {
		address = c.Lis.Addr().String()
	} else {
		address = c.Address
	}

	return &Server{
		lis:        c.Lis,
		grpcServer: c.GrpcServer,
		address:    address,
	}
}

// Running returns true if the server is running.
func (s *Server) Running() bool { return s.running }

// Serve starts the grpc server.
func (s *Server) Serve(ctx context.Context) error {
	if s.running {
		return ErrServerAlreadyRunning
	}

	if s.grpcServer == nil {
		return errors.New("grpc server is not set")
	}

	s.running = true
	defer func() {
		s.running = false
	}()

	if s.lis == nil {
		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			return fmt.Errorf("failed to create listener: %w", err)
		}
		s.lis = lis
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- s.grpcServer.Serve(s.lis)
	}()

	select {
	case <-ctx.Done():
		s.grpcServer.GracefulStop()
		return ctx.Err()
	case err := <-errChan:
		if err == nil || errors.Is(err, grpc.ErrServerStopped) {
			return nil
		}

		return fmt.Errorf("failed to serve: %w", err)
	}
}
