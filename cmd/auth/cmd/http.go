package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dwethmar/lingo/cmd/config"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func runGateway(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithCancel(context.Background())

	// Set up channel to receive signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		s := <-sigs
		slog.Info("Signal received", slog.String("signal", s.String()))
		cancel()
	}()

	config := config.New()
	s, err := setupHTTPServer(ctx, config)
	if err != nil {
		return err
	}

	g := new(errgroup.Group)
	g.Go(func() error { return s.Serve(ctx) })

	return g.Wait()
}

func NewGatewayCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "auth-gateway",
		Short: "Start the auth gateway service",
		RunE:  runGateway,
	}
}
