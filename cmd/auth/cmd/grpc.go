package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dwethmar/lingo/cmd/config"
	"github.com/dwethmar/lingo/pkg/database"
	protoauth "github.com/dwethmar/lingo/proto/gen/go/public/auth/v1"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// runAuth runs the auth server.
func runAuth(_ *cobra.Command, _ []string) error {
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

	config := config.New()

	dbURL, err := config.DatabaseURL()
	if err != nil {
		return fmt.Errorf("failed to get database url: %w", err)
	}

	db, err := database.ConnectPostgres(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			logger.Error("Failed to close database", slog.String("error", err.Error()))
		}
	}()

	auth, err := setupAuth(logger, config, db)
	if err != nil {
		return fmt.Errorf("failed to setup relay app: %w", err)
	}

	authServer := setupService(auth)
	grpcServer, err := setupServer(config, []func(grpc.ServiceRegistrar){
		func(s grpc.ServiceRegistrar) { protoauth.RegisterAuthServiceServer(s, authServer) },
		func(s grpc.ServiceRegistrar) { grpc_health_v1.RegisterHealthServer(s, authServer) },
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
		Use:   "auth",
		Short: "Start the auth grpc service",
		RunE:  runAuth,
	}
}
