package user

import (
	"context"

	"github.com/extreme-business/lingo/apps/account/domain"
	"github.com/extreme-business/lingo/apps/account/storage"
	"github.com/google/uuid"
)

const (
	perPage = 25
)

type Reader struct {
	reader storage.UserReader
}

func NewReader(storage storage.UserReader) *Reader {
	return &Reader{reader: storage}
}

func (r *Reader) Get(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := r.reader.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	var u = new(domain.User)
	return u, u.FromStorage(user)
}

func (r *Reader) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.reader.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	var u = new(domain.User)
	return u, u.FromStorage(user)
}

func (r *Reader) List(ctx context.Context, page uint) ([]*domain.User, error) {
	users, err := r.reader.List(ctx, storage.Pagination{
		Limit:  perPage,
		Offset: int(page) * perPage,
	}, storage.UserOrderBy{})
	if err != nil {
		return nil, err
	}

	var out []*domain.User
	for _, user := range users {
		var u domain.User
		if err := u.FromStorage(user); err != nil {
			return nil, err
		}

		out = append(out, &u)
	}

	return out, nil
}
