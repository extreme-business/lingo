package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/extreme-business/lingo/pkg/config"
	"github.com/extreme-business/lingo/pkg/database/postgres"
	protoaccount "github.com/extreme-business/lingo/proto/gen/go/public/account/v1"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// runAccount runs the account server.
func runAccount(_ *cobra.Command, _ []string) error {
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

	db, err := postgres.Connect(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			logger.Error("Failed to close database", slog.String("error", err.Error()))
		}
	}()

	account, err := setupAccount(logger, config, db)
	if err != nil {
		return fmt.Errorf("failed to setup relay app: %w", err)
	}

	accountServer := setupService(account)
	grpcServer, err := setupServer(config, func(s grpc.ServiceRegistrar) {
		protoaccount.RegisterAccountServiceServer(s, accountServer)
		grpc_health_v1.RegisterHealthServer(s, accountServer)
	})
	if err != nil {
		return fmt.Errorf("failed to setup grpc server: %w", err)
	}

	if err = account.Init(ctx); err != nil {
		return fmt.Errorf("failed to init account app: %w", err)
	}

	g := new(errgroup.Group)
	g.Go(func() error { return grpcServer.Serve(ctx) })

	logger.Info("Waiting for servers to finish")

	return g.Wait()
}

func NewGrpcCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "account",
		Short: "Start the account grpc service",
		RunE:  runAccount,
	}
}
