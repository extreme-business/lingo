package user

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
}
