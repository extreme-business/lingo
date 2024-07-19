package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/extreme-business/lingo/apps/account/cmd"
	"github.com/extreme-business/lingo/pkg/cli"
)

func New() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	code := cli.Run(context.Background(), logger, cmd.NewGrpcCmd())
	os.Exit(code)
}

func main() { New() }
