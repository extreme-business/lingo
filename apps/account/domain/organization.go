package domain

import (
	"time"

	"github.com/extreme-business/lingo/apps/account/storage"
	"github.com/google/uuid"
)

// OrganizationCollection is the name of the organization collection.
const OrganizationCollection = "organizations"

// Organization is an organized group of users with a particular purpose.
type Organization struct {
	ID         uuid.UUID
	LegalName  string // The official name of the organization, e.g. the registered company name.
	Slug       string // The unique identifier of the organization, e.g. the subdomain.
	CreateTime time.Time
	UpdateTime time.Time
}

func (o *Organization) FromStorage(s *storage.Organization) error {
	o.ID = s.ID
	o.LegalName = s.LegalName
	o.Slug = s.Slug
	o.CreateTime = s.CreateTime
	o.UpdateTime = s.UpdateTime
	return nil
}

func (o *Organization) ToStorage(s *storage.Organization) error {
	s.ID = o.ID
	s.LegalName = o.LegalName
	s.Slug = o.Slug
	s.CreateTime = o.CreateTime
	s.UpdateTime = o.UpdateTime
	return nil
}
