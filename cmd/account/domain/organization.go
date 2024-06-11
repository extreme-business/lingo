package domain

import (
	"time"

	"github.com/google/uuid"
)

// OrganizationCollection is the name of the organization collection.
const OrganizationCollection = "organizations"

// Organization is an organized group of users with a particular purpose.
type Organization struct {
	ID         uuid.UUID
	LegalName  string // The official name of the organization, e.g. the registered company name.
	CreateTime time.Time
	UpdateTime time.Time
}
