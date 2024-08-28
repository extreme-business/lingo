package app

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user.
type User struct {
	ID          uuid.UUID
	DisplayName string
	Email       string
	Status      string
	CreateTime  time.Time
	UpdateTime  time.Time
	DeleteTime  time.Time
}
