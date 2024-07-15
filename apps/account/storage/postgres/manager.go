package postgres

import (
	"github.com/extreme-business/lingo/cmd/account/storage"
	"github.com/extreme-business/lingo/cmd/account/storage/postgres/organization"
	"github.com/extreme-business/lingo/cmd/account/storage/postgres/user"
	"github.com/extreme-business/lingo/pkg/database"
)

// NewManager creates a new manager for storage.
func NewManager(db *database.DB) *database.Manager[storage.Repositories] {
	return database.NewManager(db, func(c database.Conn) storage.Repositories {
		return storage.Repositories{
			User:         user.New(c),
			Organization: organization.New(c),
		}
	})
}
