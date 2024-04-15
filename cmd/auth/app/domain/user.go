package domain

import (
	protoauth "github.com/dwethmar/lingo/proto/gen/go/public/auth/v1"
	"github.com/google/uuid"
)

// User is a user
type User struct {
	ID            uuid.UUID
	Username      string
	Email         string
	Password      string
	Organisations []*Organisation
}

func (u *User) ToProto(in *protoauth.User) {
	in.Id = u.ID.String()
	in.Username = u.Username
	in.Email = u.Email
}
