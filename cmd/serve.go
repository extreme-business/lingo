package cmd

import (
	authcmd "github.com/dwethmar/lingo/cmd/auth/cmd"
	"github.com/spf13/cobra"

	_ "github.com/lib/pq"
)

// serveCmd represents the relay command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve lingo services",
	Long:  `serve lingo services.`,
}

func init() {
	// auth
	serveCmd.AddCommand(authcmd.NewGrpcCmd())
	serveCmd.AddCommand(authcmd.NewGatewayCmd())

	rootCmd.AddCommand(serveCmd)
}
