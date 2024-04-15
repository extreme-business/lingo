package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(context.Context, *User) (*User, error)
	Get(context.Context, uuid.UUID) (*User, error)
	Update(context.Context, *User, ...Field) (*User, error)
	Delete(context.Context, uuid.UUID) error
}

type MockRepository struct {
	CreateFunc func(context.Context, *User) (*User, error)
	GetFunc    func(context.Context, uuid.UUID) (*User, error)
	UpdateFunc func(context.Context, *User, ...Field) (*User, error)
	DeleteFunc func(context.Context, uuid.UUID) error
}

func (m *MockRepository) Create(ctx context.Context, u *User) (*User, error) {
	return m.CreateFunc(ctx, u)
}

func (m *MockRepository) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	return m.GetFunc(ctx, id)
}

func (m *MockRepository) Update(ctx context.Context, u *User, fields ...Field) (*User, error) {
	return m.UpdateFunc(ctx, u, fields...)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFunc(ctx, id)
}
