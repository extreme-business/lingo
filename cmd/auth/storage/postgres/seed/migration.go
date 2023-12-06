package seed

import (
	"github.com/dwethmar/lingo/cmd/auth/migrations"
	"github.com/dwethmar/lingo/pkg/database/dbtest"
)

// RunMigrations runs all auth migrations. It can be passed to  dbtest.SetupTestDB.
func RunMigrations(dbURL string) error { return dbtest.Migrate(dbURL, migrations.FS) }
