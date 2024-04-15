package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/dwethmar/lingo/cmd/auth/app"
	"github.com/dwethmar/lingo/cmd/auth/authentication"
	"github.com/dwethmar/lingo/cmd/auth/registration"
	"github.com/dwethmar/lingo/cmd/auth/server"
	"github.com/dwethmar/lingo/cmd/auth/storage/user/postgres"
	"github.com/dwethmar/lingo/cmd/auth/token"
	"github.com/dwethmar/lingo/cmd/config"
	"github.com/dwethmar/lingo/pkg/clock"
	"github.com/dwethmar/lingo/pkg/database"
	"github.com/dwethmar/lingo/pkg/grpcserver"
	"github.com/dwethmar/lingo/pkg/httpserver"
	"github.com/dwethmar/lingo/pkg/uuidgen"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	protoauth "github.com/dwethmar/lingo/proto/gen/go/public/auth/v1"
)

const (
	ReadTimeout     = 5 * time.Second
	WriteTimeout    = 10 * time.Second
	IdleTimeout     = 15 * time.Second
	ShutdownTimeout = 5 * time.Second
)

// setupAuth sets up the auth application.
func setupAuth(logger *slog.Logger, db database.DB) (*app.Auth, error) {
	signingKeyRegistration, err := config.SigningKeyRegistration()
	if err != nil {
		return nil, err
	}

	signingKeyAuthentication, err := config.SigningKeyAuthentication()
	if err != nil {
		return nil, err
	}

	tokenCreated := make(chan token.Created)
	go func() {
		for created := range tokenCreated {
			logger.Info("Token created", slog.String("email", created.Email), slog.String("token", created.Token))
		}
	}()

	userRepo := postgres.NewRepository(db)
	clock := clock.Default()
	uuidgen := uuidgen.Default()

	app := app.New(
		logger,
		authentication.NewManager(authentication.Config{
			Clock:                    clock,
			SigningKeyRegistration:   []byte(signingKeyRegistration),
			SigningKeyAuthentication: []byte(signingKeyAuthentication),
			UserRepo:                 userRepo,
		}),
		registration.NewManager(registration.Config{
			UUIDgen:  uuidgen,
			Clock:    clock,
			UserRepo: userRepo,
		}),
	)

	return app, nil
}

// setupRelayGrpcServer sets up a gRPC server for the relay service.
func setupService(auth *app.Auth) (*server.Service, error) {
	return server.New(auth), nil
}

// setupGrpcServer sets up a gRPC server for the relay service.
func setupServer(serverRegisters []func(*grpc.Server)) (*grpcserver.Server, error) {
	grpcPort, err := config.GRPCPort()
	if err != nil {
		return nil, err
	}

	certFile, err := config.GrpcTLSCertFile()
	if err != nil {
		return nil, err
	}

	keyFile, err := config.GrpcTLSKeyFile()
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
func setupHttpServer(ctx context.Context) (*httpserver.Server, error) {
	port, err := config.HTTPPort()
	if err != nil {
		return nil, err
	}

	authUrl, err := config.AuthUrl()
	if err != nil {
		return nil, err
	}

	certFile, err := config.HTTPTLSCertFile()
	if err != nil {
		return nil, err
	}

	keyFile, err := config.HTTPTLSKeyFile()
	if err != nil {
		return nil, err
	}

	grpcCertFile, err := config.GrpcTLSCertFile()
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
	if err := protoauth.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, authUrl, dialOptions); err != nil {
		return nil, fmt.Errorf("failed to register gateway: %w", err)
	}

	return httpserver.New(httpserver.Config{
		Addr:            fmt.Sprintf(":%d", port),
		Handler:         mux,
		ReadTimeout:     ReadTimeout,
		WriteTimeout:    WriteTimeout,
		IdleTimeout:     IdleTimeout,
		ShutdownTimeout: ShutdownTimeout,
		CertFile:        certFile,
		KeyFile:         keyFile,
		Headers:         httpserver.CorsHeaders(),
	}), nil
}
