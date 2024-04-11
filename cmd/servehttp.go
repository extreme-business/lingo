package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Start the gateway http service",
	Long:  `Start the gateway http service.`,
	RunE:  runGateway,
}

func runGateway(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())

	// Set up channel to receive signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		s := <-sigs
		slog.Info("Signal received", slog.String("signal", s.String()))
		cancel()
	}()

	s, err := setupRelayHttpServer(ctx)
	if err != nil {
		return err
	}

	g := new(errgroup.Group)
	g.Go(func() error { return s.Serve(ctx) })

	return g.Wait()
}

func init() {
	gatewayCmd.Flags().StringP("relay-url", "r", "", "address of the relay service")
	serveCmd.AddCommand(gatewayCmd)
}
