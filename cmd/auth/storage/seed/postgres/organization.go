package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dwethmar/lingo/cmd/auth/storage/organization"
)

func Organization(ctx context.Context, db *sql.Tx, u *organization.Organization) error {
	_, err := db.ExecContext(
		ctx,
		`INSERT INTO organizations (id, display_name, create_time, update_time) VALUES ($1, $2, $3, $4)`,
		u.ID,
		u.DisplayName,
		u.CreateTime,
		u.UpdateTime,
	)

	if err != nil {
		return fmt.Errorf("failed to insert organization: %w", err)
	}

	return nil
}
