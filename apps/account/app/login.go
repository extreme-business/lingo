package app

import (
	"context"

	"github.com/extreme-business/lingo/apps/account/domain"
	authentication "github.com/extreme-business/lingo/apps/account/user/authentication"
)

// LoginResult is the result of a login operation.
type LoginResult struct {
	User         *domain.User
	AccessToken  string
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
		AccessToken:  a.AccessToken,
		RefreshToken: a.RefreshToken,
	}, nil
}
