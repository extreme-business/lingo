package relay

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/dwethmar/lingo/cmd/relay/rpc"
	"github.com/dwethmar/lingo/cmd/relay/token"
	"github.com/dwethmar/lingo/gen/go/proto/relay/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// Options options for the relay server
type Options struct {
	Logger                     *slog.Logger
	Server                     *grpc.Server
	Lis                        net.Listener
	RegistrationTokenManager   *token.Manager
	AuthenticationTokenManager *token.Manager
}

// Start starts the relay server
func Start(ctx context.Context, opt *Options) error {
	if opt.Lis == nil {
		return fmt.Errorf("listener is not set in options")
	}

	if opt.Logger == nil {
		return fmt.Errorf("logger is not set in options")
	}

	if opt.RegistrationTokenManager == nil {
		return fmt.Errorf("token manager is not set in options")
	}

	errCh := make(chan error)

	go func() {
		if err := startRpcServer(opt); err != nil {
			errCh <- fmt.Errorf("failed to start rpc server: %w", err)
		}
	}()

	go func() {
		if err := startHttpServer(ctx, opt); err != nil {
			errCh <- fmt.Errorf("failed to start http server: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func startRpcServer(opt *Options) error {
	// Register the service with the server
	relay.RegisterRelayServiceServer(opt.Server, rpc.New(opt.Logger, opt.RegistrationTokenManager, opt.AuthenticationTokenManager))

	// Register reflection service on gRPC server.
	reflection.Register(opt.Server)

	// Start the server
	// Lis will be closed by the server when it is stopped.
	if err := opt.Server.Serve(opt.Lis); err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}

	return nil
}

func startHttpServer(ctx context.Context, opt *Options) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := relay.RegisterRelayServiceHandlerFromEndpoint(ctx, mux, "localhost:9090", opts)
	if err != nil {
		return fmt.Errorf("failed to register http server: %w", err)
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	if err := http.ListenAndServe(":9090", mux); err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	return nil
}
