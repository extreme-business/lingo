package user

import (
	"context"

	"github.com/extreme-business/lingo/apps/account/storage"
	"github.com/google/uuid"
)

type Repository struct {
	CreateFunc     func(context.Context, *storage.User) (*storage.User, error)
	GetFunc        func(context.Context, uuid.UUID) (*storage.User, error)
	ListFunc       func(context.Context, storage.Pagination, storage.UserOrderBy, ...storage.Condition) ([]*storage.User, error)
	GetByEmailFunc func(context.Context, string) (*storage.User, error)
	UpdateFunc     func(context.Context, *storage.User, []storage.UserField) (*storage.User, error)
	DeleteFunc     func(context.Context, uuid.UUID) error
}

func (m *Repository) Create(ctx context.Context, u *storage.User) (*storage.User, error) {
	if m.CreateFunc == nil {
		panic("CreateFunc is not implemented")
	}
	return m.CreateFunc(ctx, u)
}

func (m *Repository) Get(ctx context.Context, id uuid.UUID) (*storage.User, error) {
	if m.GetFunc == nil {
		panic("GetFunc is not implemented")
	}
	return m.GetFunc(ctx, id)
}

func (m *Repository) List(ctx context.Context, p storage.Pagination, s storage.UserOrderBy, c ...storage.Condition) ([]*storage.User, error) {
	if m.ListFunc == nil {
		panic("ListFunc is not implemented")
	}
	return m.ListFunc(ctx, p, s, c...)
}

func (m *Repository) GetByEmail(ctx context.Context, email string) (*storage.User, error) {
	if m.GetByEmailFunc == nil {
		panic("GetByEmailFunc is not implemented")
	}
	return m.GetByEmailFunc(ctx, email)
}

func (m *Repository) Update(ctx context.Context, u *storage.User, fields []storage.UserField) (*storage.User, error) {
	if m.UpdateFunc == nil {
		panic("UpdateFunc is not implemented")
	}
	return m.UpdateFunc(ctx, u, fields)
}

func (m *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc == nil {
		panic("DeleteFunc is not implemented")
	}
	return m.DeleteFunc(ctx, id)
}
