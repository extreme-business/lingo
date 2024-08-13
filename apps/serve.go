package cmd

import (
	accountcmd "github.com/extreme-business/lingo/apps/account/cmd"
	cmscmd "github.com/extreme-business/lingo/apps/cms/cmd"
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
	// serve
	serveCmd := NewServeCmd()
	// add serve subcommands
	serveCmd.AddCommand(accountcmd.NewGrpcCmd())
	serveCmd.AddCommand(accountcmd.NewGatewayCmd())
	serveCmd.AddCommand(cmscmd.NewHTMLCmd())

	rootCmd.AddCommand(serveCmd)
}
