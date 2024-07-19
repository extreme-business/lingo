package cmd

import (
	"fmt"

	"github.com/extreme-business/lingo/apps/cms/auth"
	"github.com/extreme-business/lingo/apps/cms/server"
	"github.com/extreme-business/lingo/pkg/config"
	accountproto "github.com/extreme-business/lingo/proto/gen/go/public/account/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func runCms(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	config := config.New()

	port, err := config.HTTPPort()
	if err != nil {
		return fmt.Errorf("failed to get http port: %w", err)
	}

	accountServiceAddr, err := config.AccountURL()
	if err != nil {
		return fmt.Errorf("failed to get account service address: %w", err)
	}

	var opts = []grpc.DialOption{}
	conn, err := grpc.NewClient(accountServiceAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to create account service client: %w", err)
	}
	defer conn.Close()

	accountService := accountproto.NewAccountServiceClient(conn)

	authenticator := auth.NewAuthenticator(accountService)

	server := server.New(fmt.Sprintf(":%d", port), authenticator)
	if err := server.Serve(ctx); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func NewHtmlCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cms",
		Short: "Start the cms service",
		RunE:  runCms,
	}
}
