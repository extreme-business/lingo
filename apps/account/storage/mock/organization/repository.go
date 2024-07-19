package organization

import (
	"context"

	"github.com/extreme-business/lingo/apps/account/storage"
	"github.com/google/uuid"
)

type Repository struct {
	CreateFunc func(context.Context, *storage.Organization) (*storage.Organization, error)
	GetFunc    func(context.Context, uuid.UUID) (*storage.Organization, error)
	ListFunc   func(context.Context, storage.Pagination, storage.OrganizationOrderBy, ...storage.Condition) ([]*storage.Organization, error)
	UpdateFunc func(context.Context, *storage.Organization, ...storage.OrganizationField) (*storage.Organization, error)
	DeleteFunc func(context.Context, uuid.UUID) error
}

func (m *Repository) Create(ctx context.Context, u *storage.Organization) (*storage.Organization, error) {
	if m.CreateFunc == nil {
		panic("CreateFunc is not implemented")
	}
	return m.CreateFunc(ctx, u)
}

func (m *Repository) Get(ctx context.Context, id uuid.UUID) (*storage.Organization, error) {
	if m.GetFunc == nil {
		panic("GetFunc is not implemented")
	}
	return m.GetFunc(ctx, id)
}

func (m *Repository) List(ctx context.Context, p storage.Pagination, s storage.OrganizationOrderBy, c ...storage.Condition) ([]*storage.Organization, error) {
	if m.ListFunc == nil {
		panic("ListFunc is not implemented")
	}
	return m.ListFunc(ctx, p, s, c...)
}

func (m *Repository) Update(ctx context.Context, u *storage.Organization, fields []storage.OrganizationField) (*storage.Organization, error) {
	if m.UpdateFunc == nil {
		panic("UpdateFunc is not implemented")
	}
	return m.UpdateFunc(ctx, u, fields...)
}

func (m *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc == nil {
		panic("DeleteFunc is not implemented")
	}
	return m.DeleteFunc(ctx, id)
}
