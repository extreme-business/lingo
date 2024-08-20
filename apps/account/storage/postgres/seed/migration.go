package seed

import (
	"context"
	"testing"

	"github.com/extreme-business/lingo/apps/account/migrations"
	"github.com/extreme-business/lingo/pkg/database/dbtest"
)

// RunMigrations runs all account migrations. It can be passed to  dbtest.SetupTestDB.
func RunMigrations(ctx context.Context, t *testing.T, dbURL string) error {
	t.Helper()
	return dbtest.Migrate(ctx, dbURL, migrations.FS)
}
