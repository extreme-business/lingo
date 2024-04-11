package domain

import "github.com/google/uuid"

// User is a user
type User struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
}
