package organization

import (
	"time"

	"github.com/google/uuid"
)

// Field is a organization field.
type Field string

const (
	DisplayName Field = "display_name"
	CreateTime  Field = "create_time"
	UpdateTime  Field = "update_time"
)

// Organization is a organization.
type Organization struct {
	ID          uuid.UUID
	DisplayName string
	CreateTime  time.Time
	UpdateTime  time.Time
}
