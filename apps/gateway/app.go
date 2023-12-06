package gateway

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dwethmar/lingo/protogen/go/proto/relay/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Options struct {
	Logger          *slog.Logger
	GrpcDialOptions []grpc.DialOption
	Port            int
	RelayUrl        string
}

func Start(ctx context.Context, opts *Options) error {
	if opts.Logger == nil {
		return fmt.Errorf("logger is not set in options")
	}

	if opts.Port == 0 {
		return fmt.Errorf("port is not set in options")
	}

	if opts.RelayUrl == "" {
		return fmt.Errorf("relay url is not set in options")
	}

	logger := opts.Logger
	logger.Info("Starting gateway", slog.Int("port", opts.Port), slog.String("relay_url", opts.RelayUrl))

	mux := runtime.NewServeMux()
	err := relay.RegisterRelayServiceHandlerFromEndpoint(ctx, mux, opts.RelayUrl, opts.GrpcDialOptions)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", opts.Port),
		Handler: mux,
	}

	return server.ListenAndServe()
}
