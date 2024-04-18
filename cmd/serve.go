package cmd

import (
	authcmd "github.com/dwethmar/lingo/cmd/auth/cmd"
	"github.com/spf13/cobra"
)

func NewServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "serve lingo services",
		Long:  `serve lingo services.`,
	}
}

func init() {
	// auth
	serveCmd := NewServeCmd()
	// add auth subcommands
	serveCmd.AddCommand(authcmd.NewGrpcCmd())
	serveCmd.AddCommand(authcmd.NewGatewayCmd())

	rootCmd.AddCommand(serveCmd)
}
