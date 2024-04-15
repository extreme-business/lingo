package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/dwethmar/lingo/pkg/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lingo",
	Short: "lingo is a chat application that allows you to chat with your friends.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	code := cli.Run(context.Background(), logger, rootCmd)
	os.Exit(code)
}

func init() {
	viper.SetEnvPrefix("lingo")
	viper.AutomaticEnv()
}
