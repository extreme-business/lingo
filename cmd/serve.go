package cmd

import (
	"time"

	"github.com/spf13/cobra"

	_ "github.com/lib/pq"
)

const (
	// defaultPort default port to listen on
	defaultPort     = 8080
	ReadTimeout     = 5 * time.Second
	WriteTimeout    = 10 * time.Second
	IdleTimeout     = 15 * time.Second
	ShutdownTimeout = 5 * time.Second
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
