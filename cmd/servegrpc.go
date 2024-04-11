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

	server, err := setupRelayGrpcServer(relay)
	if err != nil {
		return fmt.Errorf("failed to setup relay server: %w", err)
	}

	g := new(errgroup.Group)
	g.Go(func() error { return server.Serve(ctx) })

	logger.Info("Waiting for servers to finish")

	return g.Wait()
}

func init() {
	relayCmd.Flags().StringP("db_url", "d", "", "Database connection string")
	serveCmd.AddCommand(relayCmd)
}
