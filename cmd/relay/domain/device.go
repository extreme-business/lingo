package domain

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID         uuid.UUID
	AccountID  uuid.UUID
	Name       string
	PublicKey  []byte
	CreateTime time.Time
}
