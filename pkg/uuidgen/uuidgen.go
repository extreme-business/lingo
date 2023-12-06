package uuidgen

import "github.com/google/uuid"

// Generator generates a new UUID.
type Generator func() uuid.UUID

// Default returns a new generator that uses uuid.New.
func Default() Generator { return func() uuid.UUID { return uuid.Must(uuid.NewV7()) } }
