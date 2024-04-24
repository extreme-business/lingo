package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dwethmar/lingo/cmd/auth/app"
	"github.com/dwethmar/lingo/cmd/auth/bootstrapping"
	"github.com/dwethmar/lingo/cmd/auth/server"
	"github.com/dwethmar/lingo/cmd/auth/storage/user/postgres"
	"github.com/dwethmar/lingo/cmd/auth/token"
	"github.com/dwethmar/lingo/cmd/auth/user/authentication"
	"github.com/dwethmar/lingo/cmd/auth/user/registration"
	"github.com/dwethmar/lingo/cmd/config"
	"github.com/dwethmar/lingo/pkg/clock"
	"github.com/dwethmar/lingo/pkg/database"
	"github.com/dwethmar/lingo/pkg/grpcserver"
	"github.com/dwethmar/lingo/pkg/httpserver"
	"github.com/dwethmar/lingo/pkg/uuidgen"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	protoauth "github.com/dwethmar/lingo/proto/gen/go/public/auth/v1"
)

const (
	readTimeout     = 5 * time.Second
	writeTimeout    = 10 * time.Second
	idleTimeout     = 15 * time.Second
	shutdownTimeout = 5 * time.Second
)

// setupAuth sets up the auth application.
func setupAuth(
	logger *slog.Logger,
	config *config.Config,
	db database.DB,
) (*app.Auth, error) {
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

	userRepo := postgres.New(db)
	clock := clock.Default()
	uuidgen := uuidgen.Default()

	app := app.New(
		logger,
		bootstrapping.NewInitializer(bootstrapping.Config{
			SystemUserID:     uuidgen.New(),
			SystemUserEmail:  "system@system.nl",
			OrganizationID:   uuidgen.New(),
			OrganizationName: "system",
			UserRepo:         userRepo,
		}),
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
func setupService(auth *app.Auth) *server.Service {
	return server.New(auth)
}

// setupGrpcServer sets up a gRPC server for the relay service.
func setupServer(config *config.Config, serviceRegistrars []func(grpc.ServiceRegistrar)) (*grpcserver.Server, error) {
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
		return nil, fmt.Errorf("failed to load TLS keys: %w", err)
	}

	return grpcserver.New(
		grpcserver.WithGrpcServer(grpc.NewServer(grpc.Creds(creds))),
		grpcserver.WithAddress(fmt.Sprintf(":%d", grpcPort)),
		grpcserver.WithServiceRegistrars(serviceRegistrars),
	), nil
}

// https://github.com/youngderekm/grpc-cookies-example/blob/master/cmd/gateway/gateway.go
func gatewayMetadataAnnotator(_ context.Context, r *http.Request) metadata.MD {
	// read token from cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil
	}

	return metadata.Pairs("token", cookie.Value)
}

func gatewayResponseModifier(ctx context.Context, r http.ResponseWriter, m proto.Message) error {
	// check if login response
	if msg, ok := m.(*protoauth.LoginUserResponse); ok {
		tokenExp, err := token.ExtractExpirationTime(msg.Token)
		if err != nil {
			return err
		}

		http.SetCookie(r, &http.Cookie{
			Name:     "token",
			Value:    msg.Token,
			SameSite: http.SameSiteStrictMode,
			Expires:  tokenExp,
		})

		refreshTokenExp, err := token.ExtractExpirationTime(msg.RefreshToken)
		if err != nil {
			return err
		}

		http.SetCookie(r, &http.Cookie{
			Name:     "refresh_token",
			Value:    msg.RefreshToken,
			Expires:  refreshTokenExp,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/v1/refresh",
		})

		msg.Token = ""
		msg.RefreshToken = ""
	}

	return nil
}

// setupRelayHttpServer sets up a HTTP server for the relay service.
func setupHTTPServer(ctx context.Context, config *config.Config) (*httpserver.Server, error) {
	port, err := config.HTTPPort()
	if err != nil {
		return nil, err
	}

	authURL, err := config.AuthURL()
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
		return nil, fmt.Errorf("failed to load TLS keys: %w", err)
	}

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	mux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(gatewayResponseModifier),
		runtime.WithMetadata(gatewayMetadataAnnotator),
	)
	if err = protoauth.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, authURL, dialOptions); err != nil {
		return nil, fmt.Errorf("failed to register gateway: %w", err)
	}

	return httpserver.New(
		httpserver.WithAddr(fmt.Sprintf(":%d", port)),
		httpserver.WithHandler(mux),
		httpserver.WithReadTimeout(readTimeout),
		httpserver.WithWriteTimeout(writeTimeout),
		httpserver.WithIdleTimeout(idleTimeout),
		httpserver.WithShutdownTimeout(shutdownTimeout),
		httpserver.WithHeaders(httpserver.CorsHeaders()),
		httpserver.WithTLS(certFile, keyFile),
	), nil
}
