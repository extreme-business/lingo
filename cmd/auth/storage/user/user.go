package user

import (
	"errors"
	"time"

	"github.com/dwethmar/lingo/cmd/auth/domain"
	"github.com/google/uuid"
)

var (
	ErrNotFound         = errors.New("user not found")
	ErrNoFieldsToUpdate = errors.New("no fields to update")
	// Unique constraint errors
	ErrUniqueIDConflict       = errors.New("unique id conflict")
	ErrUniqueUsernameConflict = errors.New("unique username conflict")
	ErrUniqueEmailConflict    = errors.New("unique email conflict")
)

type Field string

const (
	Username   Field = "username"
	Email      Field = "email"
	Password   Field = "password"
	CreateTime Field = "create_time"
	UpdateTime Field = "update_time"
)

type User struct {
	ID         uuid.UUID
	Username   string
	Email      string
	Password   string
	CreateTime time.Time
	UpdateTime time.Time
}

// ToDomain maps a User to a domain.User
func (u *User) ToDomain(in *domain.User) {
	in.ID = u.ID
	in.Username = u.Username
	in.Email = u.Email
	in.Password = u.Password
	in.CreateTime = u.CreateTime
	in.UpdateTime = u.UpdateTime
}

// FromDomain maps a domain.User to a User
func (u *User) FromDomain(in *domain.User) {
	u.ID = in.ID
	u.Username = in.Username
	u.Email = in.Email
	u.Password = in.Password
	u.CreateTime = in.CreateTime
	u.UpdateTime = in.UpdateTime
}
