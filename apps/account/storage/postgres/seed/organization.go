package seed

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/extreme-business/lingo/apps/account/storage"
	"github.com/google/uuid"
)

func NewOrganization(id, legalName, slug string, createTime, updateTime time.Time) *storage.Organization {
	return &storage.Organization{
		ID:         uuid.MustParse(id),
		LegalName:  legalName,
		Slug:       slug,
		CreateTime: createTime,
		UpdateTime: updateTime,
	}
}

func InsertOrganization(ctx context.Context, db *sql.Tx, u *storage.Organization) error {
	_, err := db.ExecContext(
		ctx,
		`INSERT INTO organizations (id, legal_name, slug, create_time, update_time) VALUES ($1, $2, $3, $4, $5)`,
		u.ID,
		u.LegalName,
		u.Slug,
		u.CreateTime,
		u.UpdateTime,
	)

	if err != nil {
		return fmt.Errorf("failed to insert organization: %w", err)
	}

	return nil
}
