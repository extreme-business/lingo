package domain

import (
	"time"

	"github.com/google/uuid"
)

// Event is a message that is sent to a recipient
type Event struct {
	ID              uuid.UUID
	SenderAccountID uuid.UUID
	Body            []byte
	CreateTime      time.Time
}

// EventRecipient is a recipient of an event
type EventRecipient struct {
	EventID       uuid.UUID
	AccountID     uuid.UUID
	DeviceID      uuid.UUID
	DecryptionKey []byte
}
