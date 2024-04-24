package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dwethmar/lingo/cmd/auth/storage/user"
)

func User(ctx context.Context, db *sql.Tx, u *user.User) error {
	_, err := db.ExecContext(
		ctx,
		`INSERT INTO users (id, organization_id, display_name, email, password, create_time, update_time) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		u.ID,
		u.OrganizationID,
		u.DisplayName,
		u.Email,
		u.Password,
		u.CreateTime,
		u.UpdateTime,
	)

	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}
