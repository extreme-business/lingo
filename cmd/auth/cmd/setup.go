package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dwethmar/lingo/cmd/auth/app"
	"github.com/dwethmar/lingo/cmd/auth/bootstrapping"
	"github.com/dwethmar/lingo/cmd/auth/config"
	"github.com/dwethmar/lingo/cmd/auth/domain"
	"github.com/dwethmar/lingo/cmd/auth/server"
	"github.com/dwethmar/lingo/cmd/auth/storage/postgres"
	"github.com/dwethmar/lingo/cmd/auth/token"
	"github.com/dwethmar/lingo/cmd/auth/user/authentication"
	"github.com/dwethmar/lingo/cmd/auth/user/registration"
	"github.com/dwethmar/lingo/pkg/clock"
	"github.com/dwethmar/lingo/pkg/database"
	"github.com/dwethmar/lingo/pkg/grpcserver"
	"github.com/dwethmar/lingo/pkg/httpserver"
	"github.com/dwethmar/lingo/pkg/resource"
	"github.com/dwethmar/lingo/pkg/uuidgen"
	protoauth "github.com/dwethmar/lingo/proto/gen/go/public/auth/v1"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

const (
	readTimeout     = 5 * time.Second
	writeTimeout    = 10 * time.Second
	idleTimeout     = 15 * time.Second
	shutdownTimeout = 5 * time.Second
)

// getSystemUserConfig gets the system user configuration from the config.
func getSystemUserConfig(c *config.Config) (bootstrapping.SystemUserConfig, error) {
	sc := bootstrapping.SystemUserConfig{}
	id, err := c.SystemUserID()
	if err != nil {
		return sc, err
	}

	sc.ID, err = uuid.Parse(id)
	if err != nil {
		return sc, fmt.Errorf("failed to parse system user id: %w", err)
	}

	sc.Email, err = c.SystemUserEmail()
	if err != nil {
		return sc, err
	}

	sc.Password, err = c.SystemUserPassword()
	if err != nil {
		return sc, err
	}

	return sc, nil
}

// getSystemOrgConfig gets the system organization configuration from the config.
func getSystemOrgConfig(c *config.Config) (bootstrapping.SystemOrgConfig, error) {
	soc := bootstrapping.SystemOrgConfig{}
	id, err := c.SystemOrganizationID()
	if err != nil {
		return soc, err
	}

	soc.ID, err = uuid.Parse(id)
	if err != nil {
		return soc, fmt.Errorf("failed to parse system org id: %w", err)
	}

	soc.LegalName, err = c.SystemOrganizationLegalName()
	if err != nil {
		return soc, err
	}

	return soc, nil
}

// setupAuth sets up the auth application.
func setupAuth(
	logger *slog.Logger,
	config *config.Config,
	db *sql.DB,
) (*app.Auth, error) {
	signingKeyRegistration, err := config.SigningKeyRegistration()
	if err != nil {
		return nil, err
	}

	signingKeyAuthentication, err := config.SigningKeyAuthentication()
	if err != nil {
		return nil, err
	}

	clock := clock.Default()
	uuidgen := uuidgen.Default()
	dbManager := postgres.NewManager(database.NewDB(db))
	repos := dbManager.Op()

	suc, err := getSystemUserConfig(config)
	if err != nil {
		return nil, err
	}

	soc, err := getSystemOrgConfig(config)
	if err != nil {
		return nil, err
	}

	app := app.New(
		logger,
		bootstrapping.New(bootstrapping.Config{
			SystemUserConfig:         suc,
			SystemOrganizationConfig: soc,
			Clock:                    clock,
			DBManager:                dbManager,
		}),
		authentication.NewManager(authentication.Config{
			Clock:                    clock,
			SigningKeyRegistration:   []byte(signingKeyRegistration),
			SigningKeyAuthentication: []byte(signingKeyAuthentication),
			UserRepo:                 repos.User,
		}),
		registration.NewManager(registration.Config{
			UUIDgen:  uuidgen,
			Clock:    clock,
			UserRepo: repos.User,
		}),
	)

	return app, nil
}

// setupRelayGrpcServer sets up a gRPC server for the relay service.
func setupService(auth *app.Auth) *server.Server {
	resourceParser := resource.NewParser()
	resourceParser.RegisterChild(domain.OrganizationCollection, domain.UserCollection)
	return server.New(auth, resourceParser)
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
		grpcserver.WithReflection(),
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

func gatewayResponseModifier(_ context.Context, r http.ResponseWriter, m proto.Message) error {
	// check if login response
	if msg, ok := m.(*protoauth.LoginUserResponse); ok {
		tokenExp, err := token.ExpirationTime(msg.GetToken())
		if err != nil {
			return err
		}

		http.SetCookie(r, &http.Cookie{
			Name:     "token",
			Value:    msg.GetToken(),
			SameSite: http.SameSiteStrictMode,
			Expires:  tokenExp,
		})

		refreshTokenExp, err := token.ExpirationTime(msg.GetRefreshToken())
		if err != nil {
			return err
		}

		http.SetCookie(r, &http.Cookie{
			Name:     "refresh_token",
			Value:    msg.GetRefreshToken(),
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
