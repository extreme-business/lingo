package postgres

import (
	"github.com/dwethmar/lingo/cmd/auth/storage"
	"github.com/dwethmar/lingo/cmd/auth/storage/postgres/organization"
	"github.com/dwethmar/lingo/cmd/auth/storage/postgres/user"
	"github.com/dwethmar/lingo/pkg/database"
)

// NewManager creates a new manager for storage.
func NewManager(db *database.DB) *database.Manager[storage.Repositories] {
	return database.NewManager(
		db,
		func(c database.Conn) storage.Repositories {
			return storage.Repositories{
				User:         user.New(c),
				Organization: organization.New(c),
			}
		},
	)
}
