package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrInvalidSortField = errors.New("invalid sort field")
)

type Pagination struct {
	Limit  int
	Offset int
}

type Direction string

const (
	ASC  Direction = "ASC"
	DESC Direction = "DESC"
)

type Sort struct {
	Field     Field
	Direction Direction
}

type OrderBy []Sort

// Validate checks if the sort fields are valid.
func (o OrderBy) Validate() error {
	fields := Fields()
	for _, s := range o {
		if s.Field == "" {
			return ErrInvalidSortField
		}

		var found bool
		for _, f := range fields {
			if s.Field == f {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("unknown field %s: %w", s.Field, ErrInvalidSortField)
		}

		if s.Direction != ASC && s.Direction != DESC {
			return fmt.Errorf("unknown direction %s: %s", s.Direction, ErrInvalidSortField)
		}
	}

	return nil
}

type Reader interface {
	Get(context.Context, uuid.UUID) (*User, error)
	GetByEmail(context.Context, string) (*User, error)
	List(context.Context, Pagination, OrderBy, ...Condition) ([]*User, error)
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
	CreateFunc     func(context.Context, *User) (*User, error)
	GetFunc        func(context.Context, uuid.UUID) (*User, error)
	ListFunc       func(context.Context, Pagination, OrderBy, ...Condition) ([]*User, error)
	GetByEmailFunc func(context.Context, string) (*User, error)
	UpdateFunc     func(context.Context, *User, ...Field) (*User, error)
	DeleteFunc     func(context.Context, uuid.UUID) error
}

func (m *MockRepository) Create(ctx context.Context, u *User) (*User, error) {
	return m.CreateFunc(ctx, u)
}

func (m *MockRepository) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	return m.GetFunc(ctx, id)
}

func (m *MockRepository) List(ctx context.Context, p Pagination, s OrderBy, c ...Condition) ([]*User, error) {
	return m.ListFunc(ctx, p, s, c...)
}

func (m *MockRepository) GetByEmail(ctx context.Context, username string) (*User, error) {
	return m.GetByEmailFunc(ctx, username)
}

func (m *MockRepository) Update(ctx context.Context, u *User, fields ...Field) (*User, error) {
	return m.UpdateFunc(ctx, u, fields...)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFunc(ctx, id)
}
