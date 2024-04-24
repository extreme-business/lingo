package organization

import (
	"context"

	"github.com/google/uuid"
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

type Reader interface {
	Get(context.Context, uuid.UUID) (*Organization, error)
	GetByEmail(context.Context, string) (*Organization, error)
	List(context.Context, Pagination, []Sort, ...Organization) ([]*Organization, error)
}

type Writer interface {
	Create(context.Context, *Organization) (*Organization, error)
	Update(context.Context, *Organization, ...Field) (*Organization, error)
	Delete(context.Context, uuid.UUID) error
}

type Repository interface {
	Reader
	Writer
}

type MockRepository struct {
	CreateFunc     func(context.Context, *Organization) (*Organization, error)
	GetFunc        func(context.Context, uuid.UUID) (*Organization, error)
	ListFunc       func(context.Context, Pagination, []Sort, ...Condition) ([]*Organization, error)
	GetByEmailFunc func(context.Context, string) (*Organization, error)
	UpdateFunc     func(context.Context, *Organization, ...Field) (*Organization, error)
	DeleteFunc     func(context.Context, uuid.UUID) error
}

func (m *MockRepository) Create(ctx context.Context, o *Organization) (*Organization, error) {
	return m.CreateFunc(ctx, o)
}

func (m *MockRepository) Get(ctx context.Context, id uuid.UUID) (*Organization, error) {
	return m.GetFunc(ctx, id)
}

func (m *MockRepository) List(ctx context.Context, p Pagination, s []Sort, c ...Condition) ([]*Organization, error) {
	return m.ListFunc(ctx, p, s, c...)
}

func (m *MockRepository) Update(ctx context.Context, u *Organization, fields ...Field) (*Organization, error) {
	return m.UpdateFunc(ctx, u, fields...)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFunc(ctx, id)
}
