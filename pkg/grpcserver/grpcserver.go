package grpcserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RegisterServer func(server *grpc.Server)

type Option func(config) config

type config struct {
	Logger        *slog.Logger        // Logger
	ServerOptions []grpc.ServerOption // Server options
	Reflection    bool                // Enable reflection
}

type Server struct {
	logger *slog.Logger
	lis    net.Listener
	serv   *grpc.Server
}

// TCPListener returns a listener on the given port.
func TCPListener(port uint) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf(":%d", port))
}

func New(
	registerServer RegisterServer,
	lis net.Listener,
	opt ...Option,
) *Server {
	var c config
	for _, o := range opt {
		c = o(c)
	}

	if c.Logger == nil {
		c.Logger = slog.Default()
	}

	grpcServer := grpc.NewServer(c.ServerOptions...)

	if c.Reflection {
		reflection.Register(grpcServer)
	}

	registerServer(grpcServer)

	return &Server{
		logger: c.Logger,
		lis:    lis,
		serv:   grpcServer,
	}
}

// Serve starts the grpc server.
// It listens for os signals and context to gracefully shutdown the server.
func (s *Server) Serve(ctx context.Context) error {
	if s.lis == nil {
		return fmt.Errorf("listener is not set")
	}

	errChan := make(chan error, 1)
	go func() {
		if err := s.serv.Serve(s.lis); err != nil {
			errChan <- fmt.Errorf("failed to serve: %w", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	select {
	case <-ctx.Done():
		s.serv.GracefulStop()
		return nil
	case <-sigChan:
		s.serv.GracefulStop()
		return nil
	case err := <-errChan:
		return err
	}
}
