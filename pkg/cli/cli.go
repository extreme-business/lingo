package cli

import (
	"context"
	"log/slog"
)

// Command is a command that can be executed.
type Command interface {
	Name() string
	ExecuteContext(ctx context.Context) error
}

// Run executes the given command and returns the exit code.
func Run(ctx context.Context, logger *slog.Logger, cmd Command) int {
	logger.Info("Running", slog.Group("command",
		"name", cmd.Name(),
	))

	if err := cmd.ExecuteContext(ctx); err != nil {
		return 1
	}

	return 0
}
