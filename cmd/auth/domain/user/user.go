package user

import (
	"time"

	"github.com/google/uuid"
)

type Field string

const (
	Username   Field = "username"
	Email      Field = "email"
	Password   Field = "password"
	CreateTime Field = "create_time"
	UpdateTime Field = "update_time"
)

type ID = uuid.UUID

type User struct {
	ID         ID
	Username   string
	Email      string
	Password   string
	CreateTime time.Time
	UpdateTime time.Time
}
