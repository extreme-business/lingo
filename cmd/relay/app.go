package relay

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/dwethmar/lingo/cmd/relay/register"
	"github.com/dwethmar/lingo/cmd/relay/rpc"
	"github.com/dwethmar/lingo/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	relayProto "github.com/dwethmar/lingo/proto/v1/relay"
)

// Options options for the relay server
type Options struct {
	Logger     *slog.Logger
	Creds      credentials.TransportCredentials
	Lis        net.Listener
	Transactor *database.Transactor
	//services
	Register *register.Registrar
}

// Start starts the relay server
func Start(opt Options) error {
	if opt.Transactor == nil {
		return fmt.Errorf("transactor is not set in options")
	}

	if opt.Lis == nil {
		return fmt.Errorf("listener is not set in options")
	}

	if opt.Logger == nil {
		return fmt.Errorf("logger is not set in options")
	}

	if opt.Creds == nil {
		return fmt.Errorf("creds is not set in options")
	}

	if opt.Register == nil {
		return fmt.Errorf("register is not set in options")
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(opt.Creds),
	)

	// Register the service with the server
	relayProto.RegisterRelayServiceServer(grpcServer, rpc.New(opt.Logger, opt.Register))

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	// Start the server
	// Lis will be closed by the server when it is stopped.
	if err := grpcServer.Serve(opt.Lis); err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}

	return nil
}
