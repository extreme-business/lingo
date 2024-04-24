package user

import (
	"github.com/dwethmar/lingo/cmd/auth/storage"
)

const (
	CreateType = "user.create"
)

type Create struct {
	user User
}

func (c *Create) Type() string { return CreateType }

func NewCreate(user User) *Create {
	return &Create{user: user}
}

type CreateHandler struct{}

func (h *CreateHandler) Handle(repo Repository, w storage.Work) error {
	return nil
}
