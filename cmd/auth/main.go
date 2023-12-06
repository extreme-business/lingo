package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/dwethmar/lingo/cmd/auth/cmd"
	"github.com/dwethmar/lingo/pkg/cli"
)

func New() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	code := cli.Run(context.Background(), logger, cmd.NewGrpcCmd())
	os.Exit(code)
}
