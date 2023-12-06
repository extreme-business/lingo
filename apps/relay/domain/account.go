package domain

import (
	"time"

	"github.com/google/uuid"
)

// Account is record that allows access to the system
type Account struct {
	ID         uuid.UUID
	EmailHash  string
	CreateTime time.Time
}
