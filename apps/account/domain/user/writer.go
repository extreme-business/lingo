package user

import (
	"context"
	"slices"
	"time"

	"github.com/extreme-business/lingo/apps/account/domain"
	"github.com/extreme-business/lingo/apps/account/storage"
	"github.com/google/uuid"
)

type Writer struct {
	c  func() time.Time // c is the clock function.
	uw storage.UserWriter
}

func NewWriter(c func() time.Time, w storage.UserWriter) *Writer {
	return &Writer{
		c:  c,
		uw: w,
	}
}

func (w *Writer) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	var err error
	var s = new(storage.User)
	if err = u.ToStorage(s); err != nil {
		return nil, err
	}
	s, err = w.uw.Create(ctx, s)
	if err != nil {
		return nil, err
	}
	result := &domain.User{}
	if err = result.FromStorage(s); err != nil {
		return nil, err
	}
	return result, nil
}

// Update updates the user.
// its sets the update time to the current time before updating the user.
// fields are sorted before updating the user and deduplicated.
func (w *Writer) Update(ctx context.Context, u *domain.User, fields []storage.UserField) (*domain.User, error) {
	u.UpdateTime = w.c()
	var err error
	s := &storage.User{}
	if err = u.ToStorage(s); err != nil {
		return nil, err
	}
	fields = append(fields, storage.UserUpdateTime)
	slices.Sort(fields)
	s, err = w.uw.Update(ctx, s, slices.Compact(fields))
	if err != nil {
		return nil, err
	}
	result := &domain.User{}
	if err = result.FromStorage(s); err != nil {
		return nil, err
	}
	return result, nil
}

func (w *Writer) Delete(ctx context.Context, id uuid.UUID) error {
	return w.uw.Delete(ctx, id)
}
