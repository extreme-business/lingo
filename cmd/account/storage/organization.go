package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Organization is a organization.
type Organization struct {
	ID         uuid.UUID
	LegalName  string
	CreateTime time.Time
	UpdateTime time.Time
}

var (
	ErrOrganizationNotFound = errors.New("organization not found")
	// Fields.
	ErrUnknownOrganizationField = errors.New("unknown organization field")
	// Sort.
	ErrEmptyOrganizationSortField       = errors.New("invalid organization sort field")
	ErrInvalidOrganizationSortDirection = errors.New("invalid organization sort direction")
	// Unique constraint errors.
	ErrConflictOrganizationID        = errors.New("unique id conflict")
	ErrConflictOrganizationLegalName = errors.New("unique legal_name conflict")
	// Immutable errors.
	ErrImmutableOrganizationID         = errors.New("field id is read-only")
	ErrImmutableOrganizationCreateTime = errors.New("field create_time is read-only")
)

type OrganizationField string

const (
	OrganizationID         OrganizationField = "id"
	OrganizationLegalName  OrganizationField = "legal_name"
	OrganizationCreateTime OrganizationField = "create_time"
	OrganizationUpdateTime OrganizationField = "update_time"
)

// OrganizationFields returns all organization fields.
func OrganizationFields() []OrganizationField {
	return []OrganizationField{
		OrganizationID,
		OrganizationLegalName,
		OrganizationCreateTime,
		OrganizationUpdateTime,
	}
}

// OrganizationSort pairs a field with a direction.
type OrganizationSort struct {
	Field     OrganizationField
	Direction Direction
}

type OrganizationOrderBy []OrganizationSort

// Validate checks if the sort fields are valid.
func (o OrganizationOrderBy) Validate() error {
	fields := OrganizationFields()
	for _, s := range o {
		if s.Field == "" {
			return ErrEmptyOrganizationSortField
		}

		var found bool
		for _, f := range fields {
			if s.Field == f {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("%s: %w", s.Field, ErrUnknownOrganizationField)
		}

		if s.Direction != ASC && s.Direction != DESC {
			return fmt.Errorf("%s: %w", s.Direction, ErrInvalidOrganizationSortDirection)
		}
	}

	return nil
}

type OrganizationReader interface {
	Get(context.Context, uuid.UUID) (*Organization, error)
	List(context.Context, Pagination, OrganizationOrderBy, ...Condition) ([]*Organization, error)
}

type OrganizationWriter interface {
	Create(context.Context, *Organization) (*Organization, error)
	Update(context.Context, *Organization, []OrganizationField) (*Organization, error)
	Delete(context.Context, uuid.UUID) error
}

type OrganizationRepository interface {
	OrganizationReader
	OrganizationWriter
}

// OrganizationByLegalNameCondition is a search condition for organizations by legal name.
type OrganizationByLegalNameCondition struct {
	Wildcard  bool
	LegalName string
}

func (OrganizationByLegalNameCondition) condition() {}
