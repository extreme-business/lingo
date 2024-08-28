package authentication

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/extreme-business/lingo/apps/account/domain"
	"github.com/extreme-business/lingo/apps/account/domain/user"
	"github.com/extreme-business/lingo/apps/account/password"
	"github.com/extreme-business/lingo/pkg/token"
)

const (
	accountTokenDuration = 5 * time.Minute
	refreshTokenDuration = 30 * time.Minute
)

// Error type allows the linter to force the user to check all errors.
type Error error

var (
	ErrUserNotFound       Error = errors.New("user not found")
	ErrInvalidCredentials Error = errors.New("invalid credentials")
)

// Authenticator is responsible for authenticating users.
type Authenticator struct {
	credentialsValidator *credentialsValidator
	userReader           *user.Reader
	AccessTokenManager   *token.Manager
	RefreshTokenManager  *token.Manager
}

type Config struct {
	Clock                  func() time.Time
	SigningKeyAccessToken  []byte
	SigningKeyRefreshToken []byte
	UserReader             *user.Reader
}

func New(c Config) *Authenticator {
	return &Authenticator{
		credentialsValidator: newCredentialsValidator(),
		userReader:           c.UserReader,
		AccessTokenManager: token.NewManager(
			c.Clock,
			c.SigningKeyAccessToken,
			accountTokenDuration,
		),
		RefreshTokenManager: token.NewManager(
			c.Clock,
			c.SigningKeyRefreshToken,
			refreshTokenDuration,
		),
	}
}

type Credentials struct {
	Email    string
	Password []byte
}

// Authentication is the process of verifying whether someone is who they claim to be when accessing a system.
type Authentication struct {
	User         *domain.User
	AccessToken  string
	RefreshToken string
}

// Authenticate a user with the given credentials.
func (m *Authenticator) Authenticate(ctx context.Context, c Credentials) (*Authentication, error) {
	if err := m.credentialsValidator.Validate(c); err != nil {
		return nil, fmt.Errorf("could not validate: %w", err)
	}

	u, err := m.userReader.GetByEmail(ctx, c.Email)
	if err != nil {
		if err == user.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if err = password.Check([]byte(c.Password), []byte(u.HashedPassword)); err != nil {
		if err == password.ErrMismatchedHashAndPassword {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to check password: %w", err)
	}

	accountToken, err := m.AccessTokenManager.Create(u.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create account token: %w", err)
	}

	refreshToken, err := m.RefreshTokenManager.Create(u.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	u.HashedPassword = ""

	return &Authentication{
		User:         u,
		AccessToken:  accountToken,
		RefreshToken: refreshToken,
	}, nil
}
