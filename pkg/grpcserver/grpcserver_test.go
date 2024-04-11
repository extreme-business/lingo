package grpcserver

import (
	"context"
	"errors"
	"net"
	"reflect"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/test/bufconn"
)

type MockGrpcServer struct {
	ServeFunc func(lis net.Listener) error
}

func (m *MockGrpcServer) Serve(lis net.Listener) error {
	return m.ServeFunc(lis)
}

func (m *MockGrpcServer) GracefulStop() {}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		c := Config{}
		got := New(c)
		if got == nil {
			t.Errorf("New() = %v, want %v", got, reflect.TypeOf(&Server{}))
		}
	})

	t.Run("New should set reflection if set", func(t *testing.T) {
		var reflectionCalled bool
		New(Config{
			Reflection: true,
			reflectionFunc: func(s reflection.GRPCServer) {
				reflectionCalled = true
			},
		})

		if !reflectionCalled {
			t.Errorf("reflectionFunc not called")
		}
	})

	t.Run("New should not set reflection if not set", func(t *testing.T) {
		var reflectionCalled bool
		New(Config{
			Reflection: false,
			reflectionFunc: func(s reflection.GRPCServer) {
				reflectionCalled = true
			},
		})

		if reflectionCalled {
			t.Errorf("reflectionFunc called")
		}
	})

	t.Run("New should set reflectionFunc if not set", func(t *testing.T) {
		c := Config{
			Reflection: true,
		}
		New(c)
	})

	t.Run("New should call ServerRegisters", func(t *testing.T) {
		var registerCalled bool
		New(Config{
			ServerRegisters: []func(*grpc.Server){
				func(s *grpc.Server) {
					registerCalled = true
				},
			},
		})

		if !registerCalled {
			t.Errorf("ServerRegisters not called")
		}
	})
}

func TestServer_Serve(t *testing.T) {
	t.Run("Serve should stop if ctx is cancled", func(t *testing.T) {
		buffer := 101024 * 1024
		lis := bufconn.Listen(buffer)
		server := grpc.NewServer()

		s := &Server{
			lis:  lis,
			Serv: server,
		}
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

		s := &Server{
			lis: lis,
			Serv: &MockGrpcServer{
				ServeFunc: func(lis net.Listener) error {
					<-make(chan struct{})
					return nil
				},
			},
		}
		ctx := context.Background()

		go func() {
			_ = s.Serve(ctx)
		}()

		for {
			for !s.Running() {
				continue
			}

			if err := s.Serve(ctx); !errors.Is(err, ErrServerAlreadyRunning) {
				t.Errorf("Serve() = %v, want %v", err, ErrServerAlreadyRunning)
			}

			break
		}
	})

	t.Run("Serve should return error if listener is not set", func(t *testing.T) {
		server := grpc.NewServer()
		s := &Server{
			Serv: server,
		}
		ctx := context.Background()

		if err := s.Serve(ctx); !errors.Is(err, ErrListenerNotSet) {
			t.Errorf("Serve() = %v, want %v", err, ErrListenerNotSet)
		}
	})

	t.Run("Serve should return error if Serve() returns error", func(t *testing.T) {
		buffer := 101024 * 1024
		lis := bufconn.Listen(buffer)

		var err = errors.New("error")

		s := &Server{
			lis: lis,
			Serv: &MockGrpcServer{
				ServeFunc: func(lis net.Listener) error {
					return err
				},
			},
		}

		ctx := context.Background()

		if !errors.Is(s.Serve(ctx), err) {
			t.Errorf("Serve() = %v, want %v", s.Serve(ctx), err)
		}
	})

	t.Run("Serve should return nil if Serve() returns nil", func(t *testing.T) {
		buffer := 101024 * 1024
		lis := bufconn.Listen(buffer)

		s := &Server{
			lis: lis,
			Serv: &MockGrpcServer{
				ServeFunc: func(lis net.Listener) error {
					return nil
				},
			},
		}

		ctx := context.Background()

		if err := s.Serve(ctx); err != nil {
			t.Errorf("Serve() = %v, want %v", err, nil)
		}
	})
}
