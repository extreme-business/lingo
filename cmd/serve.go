package cmd

import (
	"github.com/dwethmar/lingo/cmd/auth"
	"github.com/dwethmar/lingo/cmd/relay"
	"github.com/spf13/cobra"

	_ "github.com/lib/pq"
)

// ServeCmd represents the relay command
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve lingo services",
	Long:  `serve lingo services.`,
}

func init() {
	rootCmd.AddCommand(ServeCmd)

	// auth
	ServeCmd.AddCommand(auth.NewGrpcCmd())

	// relay
	ServeCmd.AddCommand(relay.NewGrpcCmd())
	ServeCmd.AddCommand(relay.NewGatewayCmd())
}
