package domain

import "github.com/google/uuid"

type Organisation struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
}
