package cmd

import (
	"database/sql"
	"fmt"

	"github.com/dwethmar/lingo/cmd/relay"
	"github.com/dwethmar/lingo/database"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	db, err := sql.Open("postgres", viper.GetString("db_url"))
	if err != nil {
		return fmt.Errorf("could not open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("could not ping db: %w", err)
	}

	defer db.Close()

	return relay.Start(relay.Options{
		Transactor: database.New(db),
	})
}

func init() {
	rootCmd.AddCommand(relayCmd)
}
