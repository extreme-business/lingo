package relay

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/dwethmar/lingo/cmd/config"
	"github.com/dwethmar/lingo/pkg/database"
	protorelay "github.com/dwethmar/lingo/proto/gen/go/private/relay/v1"
)

// runRelay runs the relay server
func runRelay(_ *cobra.Command, _ []string) error {
	logger := slog.Default()

	ctx, cancel := context.WithCancel(context.Background())

	// Set up channel to receive signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		s := <-sigs
		logger.Info("Signal received", slog.String("signal", s.String()))
		cancel()
	}()

	dbUrl, err := config.DatabaseURL()
	if err != nil {
		return fmt.Errorf("failed to get database url: %w", err)
	}

	db, dbClose, err := database.Connect(dbUrl)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	defer func() {
		if err := dbClose(); err != nil {
			logger.Error("Failed to close database", slog.String("error", err.Error()))
		}
	}()

	relay, err := setupRelay(logger, db)
	if err != nil {
		return fmt.Errorf("failed to setup relay app: %w", err)
	}

	relayServer, err := setupGrpcService(relay)
	if err != nil {
		return fmt.Errorf("failed to setup relay server: %w", err)
	}

	grpcServer, err := setupGrpcServer([]func(*grpc.Server){
		func(s *grpc.Server) { protorelay.RegisterRelayServiceServer(s, relayServer) },
		func(s *grpc.Server) { grpc_health_v1.RegisterHealthServer(s, relayServer) },
	})
	if err != nil {
		return fmt.Errorf("failed to setup grpc server: %w", err)
	}

	g := new(errgroup.Group)
	g.Go(func() error { return grpcServer.Serve(ctx) })

	logger.Info("Waiting for servers to finish")

	return g.Wait()
}

func NewGrpcCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "relay",
		Short: "Start the relay service",
		RunE:  runRelay,
	}
}
