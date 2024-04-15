package cmd

import (
	"github.com/dwethmar/lingo/cmd/auth"
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
	// auth
	ServeCmd.AddCommand(auth.NewGrpcCmd())
	ServeCmd.AddCommand(auth.NewGatewayCmd())

	rootCmd.AddCommand(ServeCmd)
}
