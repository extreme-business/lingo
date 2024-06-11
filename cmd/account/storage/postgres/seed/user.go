package seed

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dwethmar/lingo/cmd/account/storage"
	"github.com/google/uuid"
)

func NewUser(id string, organizationID string, displayName string, email string, password string, createTime time.Time, updateTime time.Time) *storage.User {
	return &storage.User{
		ID:             uuid.MustParse(id),
		OrganizationID: uuid.MustParse(organizationID),
		DisplayName:    displayName,
		Email:          email,
		Password:       password,
		CreateTime:     createTime,
		UpdateTime:     updateTime,
	}
}

func InsertUser(ctx context.Context, db *sql.Tx, u *storage.User) error {
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
