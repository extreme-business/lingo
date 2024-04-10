package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dwethmar/lingo/apps/relay"
	"github.com/dwethmar/lingo/apps/relay/server"
	"github.com/dwethmar/lingo/apps/relay/token"
	"github.com/dwethmar/lingo/pkg/clock"
	"github.com/dwethmar/lingo/pkg/grpcserver"
	"github.com/dwethmar/lingo/pkg/httpserver"
	protorelay "github.com/dwethmar/lingo/protogen/go/proto/relay/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

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

// relayCmd represents the relay command for rpc
var relayCmd = &cobra.Command{
	Use:   "relay",
	Short: "Start the relay server rpc service",
	Long:  `Start the relay server rpc service.`,
	RunE:  runRelay,
}

// runRelay runs the relay server
func runRelay(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
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

	signingKeyRegistration := viper.GetString("SIGNING_KEY_REGISTRATION")
	if signingKeyRegistration == "" {
		return fmt.Errorf("SIGNING_KEY_REGISTRATION is not set")
	}

	signingKeyAuthentication := viper.GetString("SIGNING_KEY_AUTHENTICATION")
	if signingKeyAuthentication == "" {
		return fmt.Errorf("SIGNING_KEY_AUTHENTICATION is not set")
	}

	tokenCreated := make(chan token.Created)
	go func() {
		for created := range tokenCreated {
			logger.Info("Token created", slog.String("email", created.Email), slog.String("token", created.Token))
		}
	}()

	clock := clock.New(time.UTC)

	relay := relay.New(
		logger,
		token.NewManager(
			clock,
			[]byte(signingKeyRegistration),
			15*time.Minute,
			tokenCreated,
		),
		token.NewManager(
			clock,
			[]byte(signingKeyAuthentication),
			5*time.Minute,
			tokenCreated,
		),
	)

	grpcPort := viper.GetInt("grpc_port")
	grpcAddress := fmt.Sprintf(":%d", grpcPort)

	httpPort := viper.GetInt("http_port")
	httpAddress := fmt.Sprintf(":%d", httpPort)

	certFile := viper.GetString("tls_cert_file")
	keyFile := viper.GetString("tls_key_file")

	ctx, cancel := context.WithCancel(context.Background())

	// Set up channel to receive signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		s := <-sigs
		logger.Info("Signal received", slog.String("signal", s.String()))
		cancel()
	}()

	g := new(errgroup.Group)

	// start the grpc server
	g.Go(func() error {
		logger.Info("Starting grpc server", slog.String("address", grpcAddress))

		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			return fmt.Errorf("failed to load TLS keys: %v", err)
		}

		lis, err := net.Listen("tcp", grpcAddress)
		if err != nil {
			return fmt.Errorf("failed to create listener: %w", err)
		}

		server := grpcserver.New(grpcserver.Config{
			Logger:   logger,
			Listener: lis,
			ServerOptions: []grpc.ServerOption{
				grpc.Creds(creds),
			},
			ServerRegisters: []func(*grpc.Server){
				func(s *grpc.Server) { protorelay.RegisterRelayServiceServer(s, server.New(relay)) },
			},
			Reflection: true,
		})

		if err := server.Serve(ctx); err != nil {
			logger.Error("error", slog.String("error", err.Error()))
			return fmt.Errorf("failed to serve: %w", err)
		}

		return nil
	})

	// start the http gateway
	g.Go(func() error {
		logger.Info("Starting http gateway", slog.String("address", httpAddress))

		creds, err := credentials.NewClientTLSFromFile(certFile, "lingo")
		if err != nil {
			return fmt.Errorf("failed to load TLS keys: %v", err)
		}

		mux := runtime.NewServeMux()
		if err := protorelay.RegisterRelayServiceHandlerFromEndpoint(ctx, mux, grpcAddress, []grpc.DialOption{
			grpc.WithTransportCredentials(creds),
		}); err != nil {
			return fmt.Errorf("failed to register gateway: %w", err)
		}

		server := httpserver.New(httpserver.Config{
			Addr:            httpAddress,
			Handler:         mux,
			ReadTimeout:     ReadTimeout,
			WriteTimeout:    WriteTimeout,
			IdleTimeout:     IdleTimeout,
			ShutdownTimeout: ShutdownTimeout,
			CertFile:        certFile,
			KeyFile:         keyFile,
		})

		if err := server.Serve(ctx); err != nil {
			return fmt.Errorf("failed to serve: %w", err)
		}

		return nil
	})

	logger.Info("Waiting for servers to finish")

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error: %w", err)
	}

	return nil
}

func setupEnv() error {
	if err := viper.BindEnv(
		"DB_URL",
		"GRPC_PORT",
		"HTTP_PORT",
		"TLS_CERT_FILE",
		"TLS_KEY_FILE",
		"RELAY_URL",
		"SIGNING_KEY_REGISTRATION",
		"SIGNING_KEY_AUTHENTICATION",
	); err != nil {
		return fmt.Errorf("could not bind db_url: %w", err)
	}

	if err := viper.BindPFlags(serveCmd.Flags()); err != nil {
		return fmt.Errorf("could not bind flags: %w", err)
	}

	return nil
}

func init() {
	// serve flags
	serveCmd.Flags().IntP("port", "p", defaultPort, "Port to listen on")

	// relay flags
	relayCmd.Flags().StringP("db_url", "d", "", "Database connection string")

	if err := setupEnv(); err != nil {
		panic(err)
	}

	serveCmd.AddCommand(relayCmd)
	rootCmd.AddCommand(serveCmd)
}
