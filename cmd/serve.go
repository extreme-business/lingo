package cmd

import (
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
	rootCmd.AddCommand(serveCmd)
}
