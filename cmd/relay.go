package cmd

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net"

	"github.com/dwethmar/lingo/aeslib"
	"github.com/dwethmar/lingo/cmd/relay"
	"github.com/dwethmar/lingo/cmd/relay/register"
	"github.com/dwethmar/lingo/database"
	"google.golang.org/grpc/credentials"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
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
	logger := slog.Default()

	dbConn := viper.GetString("db_url")
	if dbConn == "" {
		return fmt.Errorf("db_url is not set")
	}

	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return fmt.Errorf("could not open db: %w", err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("could not ping db: %w", err)
	}

	port := viper.GetInt("port")
	if port == 0 {
		return fmt.Errorf("port is not set")
	}

	// create a listener on TCP
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	certFile := viper.GetString("tls_cert_file")
	keyFile := viper.GetString("tls_key_file")

	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		logger.Error("failed to load TLS keys", err)

		return fmt.Errorf("failed to load TLS keys: %v", err)
	}

	keyStr := viper.GetString("AES_256_KEY")
	if keyStr == "" {
		return fmt.Errorf("AES_256_KEY is not set")
	}

	key, err := hex.DecodeString(keyStr)
	if err != nil {
		return fmt.Errorf("could not decode AES_256_KEY: %w", err)
	}

	if len(key) != 32 {
		return fmt.Errorf("key length must be 32 bytes (256 bits), got %d bytes", len(key))
	}

	logger.Info("Starting relay server", slog.Int("port", port))

	if err := relay.Start(relay.Options{
		Transactor: database.New(db),
		Lis:        lis,
		Creds:      creds,
		Register:   register.New(aeslib.New([]byte(key)), &register.LogRegisterHandler{Logger: logger}),
		Logger:     logger,
	}); err != nil {
		return fmt.Errorf("could not start relay server: %w", err)
	}

	return nil
}

func setupRelay() error {
	relayCmd.Flags().StringP("db_url", "d", "", "Database connection string")
	relayCmd.Flags().IntP("port", "p", 0, "Port to listen on")

	if err := viper.BindEnv("DB_URL"); err != nil {
		return fmt.Errorf("could not bind db_url: %w", err)
	}

	if err := viper.BindEnv("PORT"); err != nil {
		return fmt.Errorf("could not bind port: %w", err)
	}

	if err := viper.BindEnv("TLS_CERT_FILE"); err != nil {
		return fmt.Errorf("could not bind tls_cert_file: %w", err)
	}

	if err := viper.BindEnv("TLS_KEY_FILE"); err != nil {
		return fmt.Errorf("could not bind tls_key_file: %w", err)
	}

	if err := viper.BindPFlags(relayCmd.Flags()); err != nil {
		return fmt.Errorf("could not bind flags: %w", err)
	}

	return nil
}

func init() {
	if err := setupRelay(); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(relayCmd)
}
