package grpcserver_test

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/dwethmar/lingo/pkg/grpcserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type MockGrpcServer struct {
	ServeFunc          func(lis net.Listener) error
	GetServiceInfoFunc func() map[string]grpc.ServiceInfo
}

func (m *MockGrpcServer) Serve(lis net.Listener) error {
	if m.ServeFunc == nil {
		panic("ServeFunc is not set")
	}

	return m.ServeFunc(lis)
}
func (m *MockGrpcServer) RegisterService(_ *grpc.ServiceDesc, _ any) {}
func (m *MockGrpcServer) GracefulStop()                              {}
func (m *MockGrpcServer) GetServiceInfo() map[string]grpc.ServiceInfo {
	if m.GetServiceInfoFunc == nil {
		panic("GetServiceInfoFunc is not set")
	}

	return m.GetServiceInfoFunc()
}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		buffer := 101024 * 1024
		lis := bufconn.Listen(buffer)
		var visited bool

		s := grpcserver.New(
			grpcserver.WithGrpcServer(&MockGrpcServer{}),
			grpcserver.WithListener(lis),
			grpcserver.WithServiceRegistrar(func(grpc.ServiceRegistrar) {
				visited = true
			}),
			grpcserver.WithReflection(),
		)

		if s == nil {
			t.Error("expected a new server")
		}

		if !visited {
			t.Error("expected the server to be visited")
		}
	})
}

func TestServer_Serve(t *testing.T) {
	t.Run("Serve should stop if ctx is canceled", func(t *testing.T) {
		buffer := 101024 * 1024
		lis := bufconn.Listen(buffer)

		s := grpcserver.New(
			grpcserver.WithListener(lis),
			grpcserver.WithServiceRegistrar(func(grpc.ServiceRegistrar) {}),
			grpcserver.WithReflection(),
		)

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		cancel()

		if err := s.Serve(ctx); !errors.Is(err, context.Canceled) {
			t.Errorf("Serve() = %v, want %v", err, context.Canceled)
		}
	})

	t.Run("Serve should return error if server is already running", func(t *testing.T) {
		buffer := 101024 * 1024
		lis := bufconn.Listen(buffer)

		s := grpcserver.New(
			grpcserver.WithListener(lis),
			grpcserver.WithServiceRegistrar(func(grpc.ServiceRegistrar) {}),
			grpcserver.WithReflection(),
		)
		ctx := context.Background()

		go func() {
			_ = s.Serve(ctx)
		}()

		for {
			for !s.Running() {
				continue
			}

			if err := s.Serve(ctx); !errors.Is(err, grpcserver.ErrServerAlreadyRunning) {
				t.Errorf("Serve() = %v, want %v", err, grpcserver.ErrServerAlreadyRunning)
			}

			break
		}
	})

	t.Run("Serve should return error if Serve() returns error", func(t *testing.T) {
		var expectedErr = errors.New("expected error")

		s := grpcserver.New(
			grpcserver.WithGrpcServer(&MockGrpcServer{
				ServeFunc: func(_ net.Listener) error { return expectedErr },
			}),
			grpcserver.WithServiceRegistrar(func(grpc.ServiceRegistrar) {}),
			grpcserver.WithReflection(),
		)

		ctx := context.Background()
		err := s.Serve(ctx)

		if err == nil || !errors.Is(err, expectedErr) {
			t.Errorf("Serve() = %v, want %v", err, expectedErr)
		}
	})

	t.Run("Serve should return nil if server is stopped", func(t *testing.T) {
		s := grpcserver.New(
			grpcserver.WithGrpcServer(&MockGrpcServer{
				ServeFunc: func(_ net.Listener) error { return grpc.ErrServerStopped },
			}),
			grpcserver.WithServiceRegistrar(func(grpc.ServiceRegistrar) {}),
			grpcserver.WithReflection(),
		)

		ctx := context.Background()
		err := s.Serve(ctx)

		if err != nil {
			t.Errorf("Serve() = %v, want nil", err)
		}
	})
}
