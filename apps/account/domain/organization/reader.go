package organization

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
	reader storage.OrganizationReader
}

func NewReader(storage storage.OrganizationReader) *Reader {
	return &Reader{reader: storage}
}

func (r *Reader) Get(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	organization, err := r.reader.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	var o = new(domain.Organization)
	return o, o.FromStorage(organization)
}

func (r *Reader) List(ctx context.Context, page uint) ([]*domain.Organization, error) {
	organizations, err := r.reader.List(ctx, storage.Pagination{
		Limit:  perPage,
		Offset: int(page) * perPage,
	}, storage.OrganizationOrderBy{})
	if err != nil {
		return nil, err
	}
	var out []*domain.Organization
	for _, organization := range organizations {
		var o domain.Organization
		if err = o.FromStorage(organization); err != nil {
			return nil, err
		}
		out = append(out, &o)
	}
	return out, nil
}
