package relay

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/dwethmar/lingo/cmd/relay/rpc"
	"github.com/dwethmar/lingo/cmd/relay/token"
	"github.com/dwethmar/lingo/gen/go/proto/relay/v1"

	"google.golang.org/grpc"
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
func Start(ctx context.Context, opts *Options) error {
	if opts.Lis == nil {
		return fmt.Errorf("listener is not set in options")
	}

	if opts.Logger == nil {
		return fmt.Errorf("logger is not set in options")
	}

	if opts.RegistrationTokenManager == nil {
		return fmt.Errorf("registration token manager is not set in options")
	}

	if opts.AuthenticationTokenManager == nil {
		return fmt.Errorf("authentication token manager is not set in options")
	}

	// Register the service with the server
	relay.RegisterRelayServiceServer(opts.Server, rpc.New(
		opts.Logger,
		opts.RegistrationTokenManager,
		opts.AuthenticationTokenManager,
	))

	// Register reflection service on gRPC server.
	reflection.Register(opts.Server)

	// Use a channel to communicate server errors
	errChan := make(chan error, 1)
	go func() {
		if err := opts.Server.Serve(opts.Lis); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		opts.Logger.Info("Shutting down relay server")
		opts.Server.GracefulStop()
		return nil
	case err := <-errChan:
		return err
	}
}
