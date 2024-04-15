package cli

import (
	"context"
	"log/slog"

	"github.com/spf13/cobra"
)

// Run executes the given command and returns the exit code.
func Run(ctx context.Context, logger *slog.Logger, cmd *cobra.Command) int {
	logger.Info("Running", slog.Group("command",
		"name", cmd.Name(),
	))

	if err := cmd.ExecuteContext(ctx); err != nil {
		return 1
	}

	return 0
}
