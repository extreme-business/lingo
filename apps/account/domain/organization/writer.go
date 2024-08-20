package organization

import (
	"context"
	"time"

	"github.com/extreme-business/lingo/apps/account/domain"
	"github.com/extreme-business/lingo/apps/account/storage"
	"github.com/google/uuid"
)

type Writer struct {
	c  func() time.Time // c is the clock function.
	uw storage.OrganizationWriter
}

func NewWriter(c func() time.Time, w storage.OrganizationWriter) *Writer {
	return &Writer{
		c:  c,
		uw: w,
	}
}

func (w *Writer) Create(ctx context.Context, o *domain.Organization) (*domain.Organization, error) {
	var err error
	var s = new(storage.Organization)
	if err = o.ToStorage(s); err != nil {
		return nil, err
	}
	s, err = w.uw.Create(ctx, s)
	if err != nil {
		return nil, err
	}
	result := &domain.Organization{}
	if err = result.FromStorage(s); err != nil {
		return nil, err
	}
	return result, nil
}

func (w *Writer) Update(ctx context.Context, o *domain.Organization) (*domain.Organization, error) {
	o.UpdateTime = w.c()
	var err error
	s := &storage.Organization{}
	if err = o.ToStorage(s); err != nil {
		return nil, err
	}
	s, err = w.uw.Update(ctx, s, []storage.OrganizationField{
		storage.OrganizationLegalName,
	})
	if err != nil {
		return nil, err
	}
	result := &domain.Organization{}
	if err = result.FromStorage(s); err != nil {
		return nil, err
	}
	return result, nil
}

func (w *Writer) Delete(ctx context.Context, id uuid.UUID) error {
	return w.uw.Delete(ctx, id)
}
