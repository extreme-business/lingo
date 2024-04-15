package user

import (
	"time"

	"github.com/dwethmar/lingo/cmd/auth/domain"
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

type User struct {
	ID         uuid.UUID
	Username   string
	Email      string
	Password   string
	CreateTime time.Time
	UpdateTime time.Time
}

// ToDomain maps a User to a domain.User
func (u *User) ToDomain(in *domain.User) error {
	in.ID = u.ID
	in.Username = u.Username
	in.Email = u.Email
	in.Password = u.Password
	in.CreateTime = u.CreateTime
	in.UpdateTime = u.UpdateTime
	return nil
}

// FromDomain maps a domain.User to a User
func (u *User) FromDomain(in *domain.User) error {
	u.ID = in.ID
	u.Username = in.Username
	u.Email = in.Email
	return nil
}
