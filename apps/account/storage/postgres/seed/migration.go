package seed

import (
	"github.com/extreme-business/lingo/cmd/account/migrations"
	"github.com/extreme-business/lingo/pkg/database/dbtest"
)

// RunMigrations runs all account migrations. It can be passed to  dbtest.SetupTestDB.
func RunMigrations(dbURL string) error { return dbtest.Migrate(dbURL, migrations.FS) }
