package cmd

import (
	"github.com/dwethmar/lingo/cmd/relay"

	"github.com/spf13/cobra"
)

// relayCmd represents the relay command
var relayCmd = &cobra.Command{
	Use:   "relay",
	Short: "Start the relay server",
	Long: `Start the relay server. This server is responsible for
	receiving messages from the client and forwarding them to the
	appropriate client.`,
	RunE: run,
}

func run(cmd *cobra.Command, args []string) error {
	return relay.Start()
}

func init() {
	rootCmd.AddCommand(relayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// relayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// relayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
