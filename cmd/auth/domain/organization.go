package domain

import "github.com/google/uuid"

type Organization struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
}
