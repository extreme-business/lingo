package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/dwethmar/lingo/apps/relay"
	"github.com/dwethmar/lingo/apps/relay/server"
	"github.com/dwethmar/lingo/apps/relay/token"
	"github.com/dwethmar/lingo/pkg/clock"
	"github.com/dwethmar/lingo/pkg/database"
	"github.com/dwethmar/lingo/pkg/grpcserver"
	"github.com/dwethmar/lingo/pkg/httpserver"
	protorelay "github.com/dwethmar/lingo/protogen/go/proto/private/relay/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var ( // env keys
	EnvKeyDatabaseURL              = "DB_URL"
	EnvKeySigningKeyRegistration   = "SIGNING_KEY_REGISTRATION"
	EnvKeySigningKeyAuthentication = "SIGNING_KEY_AUTHENTICATION"
	EnvKeyHTTPPort                 = "HTTP_PORT"
	EnvKeyGRPCPort                 = "GRPC_PORT"
	EnvKeyGrpcTLSCertFile          = "GRPC_TLS_CERT_FILE"
	EnvKeyGrpcTLSKeyFile           = "GRPC_TLS_KEY_FILE"
	EnvKeyHTTPTLSKeyFile           = "HTTP_TLS_KEY_FILE"
	EnvKeyHTTPTLSCertFile          = "HTTP_TLS_CERT_FILE"
	EnvKeyRelayUrl                 = "RELAY_URL"
)

// getConfigString returns the value of the key as a string.
func getConfigString(key string) (string, error) {
	if !viper.IsSet(key) {
		return "", fmt.Errorf("%s is not set", key)
	}

	value := viper.GetString(key)
	if value == "" {
		return "", fmt.Errorf("%s is empty", key)
	}

	return value, nil
}

// getConfigInt returns the value of the key as an integer.
func getConfigInt(key string) (int, error) {
	if !viper.IsSet(key) {
		return 0, fmt.Errorf("%s is not set", key)
	}

	return viper.GetInt(key), nil
}

// setupDatabase sets up the database connection.
func setupDatabase() (database.DB, func() error, error) {
	dbConn, err := getConfigString(EnvKeyDatabaseURL)
	if err != nil {
		return nil, nil, err
	}

	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("could not ping db: %w", err)
	}

	return db, db.Close, nil
}

// setupRelayApp sets up the relay application.
func setupRelayApp(logger *slog.Logger, _ database.DB) (*relay.Relay, error) {
	signingKeyRegistration, err := getConfigString(EnvKeySigningKeyRegistration)
	if err != nil {
		return nil, err
	}

	signingKeyAuthentication, err := getConfigString(EnvKeySigningKeyAuthentication)
	if err != nil {
		return nil, err
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

	return relay, nil
}

// setupRelayGrpcServer sets up a gRPC server for the relay service.
func setupRelayGrpcServer(relay *relay.Relay) (*server.Server, error) {
	return server.New(relay), nil
}

// setupGrpcServer sets up a gRPC server for the relay service.
func setupGrpcServer(serverRegisters []func(*grpc.Server)) (*grpcserver.Server, error) {
	grpcPort, err := getConfigInt(EnvKeyGRPCPort)
	if err != nil {
		return nil, err
	}

	certFile, err := getConfigString(EnvKeyGrpcTLSCertFile)
	if err != nil {
		return nil, err
	}

	keyFile, err := getConfigString(EnvKeyGrpcTLSKeyFile)
	if err != nil {
		return nil, err
	}

	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS keys: %v", err)
	}

	grpcAddress := fmt.Sprintf(":%d", grpcPort)
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	return grpcserver.New(grpcserver.Config{
		Listener: lis,
		ServerOptions: []grpc.ServerOption{
			grpc.Creds(creds),
		},
		ServerRegisters: serverRegisters,
		Reflection:      true,
	}), nil
}

// setupRelayHttpServer
func setupRelayHttpServer(ctx context.Context) (*httpserver.Server, error) {
	port, err := getConfigInt(EnvKeyHTTPPort)
	if err != nil {
		return nil, err
	}

	relayUrl, err := getConfigString(EnvKeyRelayUrl)
	if err != nil {
		return nil, err
	}

	certFile, err := getConfigString(EnvKeyHTTPTLSCertFile)
	if err != nil {
		return nil, err
	}

	keyFile, err := getConfigString(EnvKeyHTTPTLSKeyFile)
	if err != nil {
		return nil, err
	}

	grpcCertFile, err := getConfigString(EnvKeyGrpcTLSCertFile)
	if err != nil {
		return nil, err
	}

	creds, err := credentials.NewClientTLSFromFile(grpcCertFile, "lingo")
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS keys: %v", err)
	}

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	mux := runtime.NewServeMux()
	if err := protorelay.RegisterRelayServiceHandlerFromEndpoint(ctx, mux, relayUrl, dialOptions); err != nil {
		return nil, fmt.Errorf("failed to register gateway: %w", err)
	}

	return httpserver.New(httpserver.Config{
		Addr:     fmt.Sprintf(":%d", port),
		Handler:  mux,
		CertFile: certFile,
		KeyFile:  keyFile,
		Cors:     true,
	}), nil
}
