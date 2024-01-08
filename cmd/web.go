package cmd

import (
	"fmt"

	"github.com/dwethmar/lingo/cmd/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Serve web interface for lingo",
	Long:  `Serve web interface for lingo. This interface is responsible for managing accounts`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("web called")
	},
}

func runWeb(cmd *cobra.Command, args []string) error {
	if err := web.Start(&web.Options{
		Port: viper.GetInt("port"),
	}); err != nil {
		return fmt.Errorf("failed to start web server: %w", err)
	}

	return nil
}

func setupWeb() error {
	webCmd.Flags().IntP("port", "p", 0, "Port to listen on")

	if err := viper.BindEnv("PORT"); err != nil {
		return fmt.Errorf("could not bind port: %w", err)
	}

	return nil
}

func init() {
	if err := setupWeb(); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(webCmd)
}
