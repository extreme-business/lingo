package cmd

import (
	accountcmd "github.com/extreme-business/lingo/cmd/account/cmd"
	"github.com/spf13/cobra"
)

func NewServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "serve lingo services",
		Long:  `serve lingo services.`,
	}
}

//nolint:gochecknoinits // This is the entry point of the serve cli.
func init() {
	// account
	serveCmd := NewServeCmd()
	// add account subcommands
	serveCmd.AddCommand(accountcmd.NewGrpcCmd())
	serveCmd.AddCommand(accountcmd.NewGatewayCmd())

	rootCmd.AddCommand(serveCmd)
}
