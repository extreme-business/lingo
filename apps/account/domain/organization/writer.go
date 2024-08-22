package organization

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

// Update updates the organization.
// its sets the update time to the current time before updating the organization.
// fields are sorted before updating the organization and deduplicated.
func (w *Writer) Update(ctx context.Context, o *domain.Organization, fields []storage.OrganizationField) (*domain.Organization, error) {
	o.UpdateTime = w.c()
	var err error
	s := &storage.Organization{}
	if err = o.ToStorage(s); err != nil {
		return nil, err
	}
	fields = append(fields, storage.OrganizationUpdateTime)
	slices.Sort(fields)
	s, err = w.uw.Update(ctx, s, slices.Compact(fields))
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
