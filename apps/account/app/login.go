package app

import (
	"context"
	"errors"

	"github.com/extreme-business/lingo/apps/account/auth/authentication"
	"github.com/extreme-business/lingo/apps/account/domain"
)

var (
	// ErrUserNotFound is returned when the user is not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidCredentials is returned when the credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// LoginResult is the result of a login operation.
type LoginResult struct {
	User         *domain.User
	AccessToken  string
	RefreshToken string
}

// LoginUser logs in a user with the given email and password.
func (r *App) LoginUser(ctx context.Context, email, password string) (*LoginResult, error) {
	a, err := r.authenticator.Authenticate(ctx, authentication.Credentials{
		Email:    email,
		Password: []byte(password),
	})
	if err != nil {
		switch err {
		case authentication.ErrInvalidCredentials:
			return nil, ErrInvalidCredentials
		case authentication.ErrUserNotFound:
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}
	return &LoginResult{
		User:         a.User,
		AccessToken:  a.AccessToken,
		RefreshToken: a.RefreshToken,
	}, nil
}
