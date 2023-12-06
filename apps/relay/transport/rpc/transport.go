package rpc

import (
	"context"
	"fmt"

	"github.com/dwethmar/lingo/apps/relay"
	"github.com/dwethmar/lingo/pkg/grpcserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	protorelay "github.com/dwethmar/lingo/protogen/go/proto/relay/v1"
)

type Transport struct {
	server *grpcserver.Server
}

type Config struct {
	Relay    *relay.Relay
	Port     uint
	CertFile string
	KeyFile  string
}

func New(c Config) (*Transport, error) {
	if c.Relay == nil {
		return nil, fmt.Errorf("relay is not set in config")
	}

	if c.Port == 0 {
		return nil, fmt.Errorf("port is not set in config")
	}

	lis, err := grpcserver.TCPListener(c.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	creds, err := credentials.NewServerTLSFromFile(c.CertFile, c.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS keys: %v", err)
	}

	s := grpcserver.New(
		func(server *grpc.Server) {
			protorelay.RegisterRelayServiceServer(server, &Server{
				relay: c.Relay,
			})
		},
		lis,
		grpcserver.WithReflection(),
		grpcserver.WithCredentials(creds),
	)

	return &Transport{
		server: s,
	}, nil
}

func (t *Transport) Serve(ctx context.Context) error {
	return t.server.Serve(ctx)
}
