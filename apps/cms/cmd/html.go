package cmd

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/extreme-business/lingo/apps/cms/app"
	"github.com/extreme-business/lingo/apps/cms/server"
	"github.com/extreme-business/lingo/apps/cms/server/token"
	"github.com/extreme-business/lingo/pkg/config"
	"github.com/extreme-business/lingo/pkg/httpmiddleware"
	accountproto "github.com/extreme-business/lingo/proto/gen/go/public/account/v1"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func runCms(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	logger := slog.Default()
	config := config.New()

	port, err := config.HTTPPort()
	if err != nil {
		return fmt.Errorf("failed to get http port: %w", err)
	}

	signingKeyAccessToken, err := config.SigningKeyAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get signing key: %w", err)
	}

	accountServiceAddr, err := config.AccountServiceURL()
	if err != nil {
		return fmt.Errorf("failed to get account service address: %w", err)
	}

	accountServiceCertFile, err := config.AccountServiceTLSCertFile()
	if err != nil {
		return fmt.Errorf("failed to get account service cert file: %w", err)
	}

	creds, err := credentials.NewClientTLSFromFile(accountServiceCertFile, "lingo")
	if err != nil {
		return fmt.Errorf("failed to load TLS keys: %w", err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	accountClient, err := grpc.NewClient(accountServiceAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to create account service client: %w", err)
	}
	defer accountClient.Close()

	accountService := accountproto.NewAccountServiceClient(accountClient)
	app := app.New(accountService)

	tokenValidator := token.NewTokenValidator([]byte(signingKeyAccessToken))
	authMiddleware := httpmiddleware.AuthCookie("access_token", tokenValidator, "/login", map[string][]string{
		"/login":    {http.MethodPost, http.MethodGet},
		"/register": {http.MethodPost, http.MethodGet},
	})

	server, err := server.New(
		logger,
		fmt.Sprintf(":%d", port),
		app,
		authMiddleware,
	)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}
	if err = server.Serve(ctx); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func NewHTMLCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cms",
		Short: "Start the cms service",
		RunE:  runCms,
	}
}
