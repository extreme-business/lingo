package postgres

import (
	"context"
	"database/sql"

	"github.com/dwethmar/lingo/cmd/auth/storage/user"
)

func User(ctx context.Context, db *sql.DB, u *user.User) error {
	_, err := db.ExecContext(
		ctx,
		`INSERT INTO users (id, username, email, password, create_time, update_time) VALUES ($1, $2, $3, $4, $5, $6)`,
		u.ID,
		u.Username,
		u.Email,
		u.Password,
		u.CreateTime,
		u.UpdateTime,
	)

	if err != nil {
		return err
	}

	return nil
}
