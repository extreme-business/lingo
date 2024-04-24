package organization

import (
	"time"

	"github.com/google/uuid"
)

// Organization is a organization.
type Organization struct {
	ID          uuid.UUID
	DisplayName string
	CreateTime  time.Time
	UpdateTime  time.Time
}
