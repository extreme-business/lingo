package app

import (
	"context"

	"github.com/dwethmar/lingo/cmd/account/domain"
	authentication "github.com/dwethmar/lingo/cmd/account/user/authentication"
)

// LoginResult is the result of a login operation.
type LoginResult struct {
	User         *domain.User
	Token        string
	RefreshToken string
}

// LoginUser logs in a user with the given email and password.
func (r *Account) LoginUser(ctx context.Context, email, password string) (*LoginResult, error) {
	a, err := r.accountManager.Authenticate(ctx, authentication.Credentials{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:         a.User,
		Token:        a.Token,
		RefreshToken: a.RefreshToken,
	}, nil
}
