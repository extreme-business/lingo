package cli

import (
	"context"
	"errors"
	"log/slog"
	"testing"
)

type mockCommand struct {
	nameFunc func() string
	execFunc func(ctx context.Context) error
}

func (m *mockCommand) Name() string                             { return m.nameFunc() }
func (m *mockCommand) ExecuteContext(ctx context.Context) error { return m.execFunc(ctx) }

func TestRun(t *testing.T) {
	t.Run("Run", func(t *testing.T) {
		t.Run("when command succeeds", func(t *testing.T) {
			logger := slog.Default()
			ctx := context.Background()

			cmd := &mockCommand{
				nameFunc: func() string { return "test" },
				execFunc: func(ctx context.Context) error { return nil },
			}

			got := Run(ctx, logger, cmd)
			if want := 0; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})

		t.Run("when command fails", func(t *testing.T) {
			logger := slog.Default()
			ctx := context.Background()

			cmd := &mockCommand{
				nameFunc: func() string { return "test" },
				execFunc: func(ctx context.Context) error { return errors.ErrUnsupported },
			}

			got := Run(ctx, logger, cmd)
			if want := 1; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	})
}
