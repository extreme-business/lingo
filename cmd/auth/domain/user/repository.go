package user

import (
	"context"
)

type Repository interface {
	Create(context.Context, *User) (*User, error)
	Get(context.Context, ID) (*User, error)
	Update(context.Context, *User, ...Field) (*User, error)
	Delete(context.Context, ID) error
}
