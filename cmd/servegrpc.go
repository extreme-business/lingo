package cmd

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

	protorelay "github.com/dwethmar/lingo/protogen/go/proto/private/relay/v1"
)

// relayCmd represents the relay command for rpc
var relayCmd = &cobra.Command{
	Use:   "relay",
	Short: "Start the relay server rpc service",
	Long:  `Start the relay server rpc service.`,
	RunE:  runRelay,
}

// runRelay runs the relay server
func runRelay(cmd *cobra.Command, args []string) error {
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

	db, dbClose, err := setupDatabase()
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}
	defer dbClose()

	relay, err := setupRelayApp(logger, db)
	if err != nil {
		return fmt.Errorf("failed to setup relay app: %w", err)
	}

	relayServer, err := setupRelayGrpcServer(relay)
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

func init() {
	relayCmd.Flags().StringP("db_url", "d", "", "Database connection string")
	serveCmd.AddCommand(relayCmd)
}
