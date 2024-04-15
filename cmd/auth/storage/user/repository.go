package user

import (
	"context"

	"github.com/google/uuid"
)

type Reader interface {
	Get(context.Context, uuid.UUID) (*User, error)
	GetByUsername(context.Context, string) (*User, error)
}

type Writer interface {
	Create(context.Context, *User) (*User, error)
	Update(context.Context, *User, ...Field) (*User, error)
	Delete(context.Context, uuid.UUID) error
}

type Repository interface {
	Reader
	Writer
}

type MockRepository struct {
	CreateFunc        func(context.Context, *User) (*User, error)
	GetFunc           func(context.Context, uuid.UUID) (*User, error)
	GetByUsernameFunc func(context.Context, string) (*User, error)
	UpdateFunc        func(context.Context, *User, ...Field) (*User, error)
	DeleteFunc        func(context.Context, uuid.UUID) error
}

func (m *MockRepository) Create(ctx context.Context, u *User) (*User, error) {
	return m.CreateFunc(ctx, u)
}

func (m *MockRepository) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	return m.GetFunc(ctx, id)
}

func (m *MockRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	return m.GetByUsernameFunc(ctx, username)
}

func (m *MockRepository) Update(ctx context.Context, u *User, fields ...Field) (*User, error) {
	return m.UpdateFunc(ctx, u, fields...)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFunc(ctx, id)
}
