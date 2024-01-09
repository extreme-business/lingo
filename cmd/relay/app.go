package relay

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/dwethmar/lingo/cmd/relay/rpc"
	"github.com/dwethmar/lingo/cmd/relay/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	relayProto "github.com/dwethmar/lingo/proto/v1/relay"
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
func Start(opt Options) error {
	if opt.Lis == nil {
		return fmt.Errorf("listener is not set in options")
	}

	if opt.Logger == nil {
		return fmt.Errorf("logger is not set in options")
	}

	if opt.RegistrationTokenManager == nil {
		return fmt.Errorf("token manager is not set in options")
	}

	// Register the service with the server
	relayProto.RegisterRelayServiceServer(opt.Server, rpc.New(opt.Logger, opt.RegistrationTokenManager, opt.AuthenticationTokenManager))

	// Register reflection service on gRPC server.
	reflection.Register(opt.Server)

	// Start the server
	// Lis will be closed by the server when it is stopped.
	if err := opt.Server.Serve(opt.Lis); err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}

	return nil
}
