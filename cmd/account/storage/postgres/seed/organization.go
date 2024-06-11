package seed

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dwethmar/lingo/cmd/account/storage"
	"github.com/google/uuid"
)

func NewOrganization(id string, legalName string, createTime time.Time, updateTime time.Time) *storage.Organization {
	return &storage.Organization{
		ID:         uuid.MustParse(id),
		LegalName:  legalName,
		CreateTime: createTime,
		UpdateTime: updateTime,
	}
}

func InsertOrganization(ctx context.Context, db *sql.Tx, u *storage.Organization) error {
	_, err := db.ExecContext(
		ctx,
		`INSERT INTO organizations (id, legal_name, create_time, update_time) VALUES ($1, $2, $3, $4)`,
		u.ID,
		u.LegalName,
		u.CreateTime,
		u.UpdateTime,
	)

	if err != nil {
		return fmt.Errorf("failed to insert organization: %w", err)
	}

	return nil
}
