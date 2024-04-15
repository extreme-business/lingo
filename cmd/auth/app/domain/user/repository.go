package user

import "context"

type Repository interface {
	Create(context.Context, *User) error
	Get(context.Context, string) (*User, error)
	Update(context.Context, *User) error
	Delete(context.Context, string) error
}
