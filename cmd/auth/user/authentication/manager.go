package authentication

import (
	"context"
	"errors"
	"time"

	"github.com/dwethmar/lingo/cmd/auth/password"
	"github.com/dwethmar/lingo/cmd/auth/storage/user"
	"github.com/dwethmar/lingo/cmd/auth/token"
	"github.com/dwethmar/lingo/pkg/clock"
)

const (
	authTokenDuration    = 5 * time.Minute
	refreshTokenDuration = 30 * time.Minute
)

type Manager struct {
	userRepo            user.Repository
	AuthTokenManager    *token.Manager
	RefreshTokenManager *token.Manager
}

type Config struct {
	Clock                    *clock.Clock
	SigningKeyRegistration   []byte
	SigningKeyAuthentication []byte
	UserRepo                 user.Repository
}

func NewManager(c Config) *Manager {
	return &Manager{
		userRepo: c.UserRepo,
		AuthTokenManager: token.NewManager(
			c.Clock,
			c.SigningKeyRegistration,
			authTokenDuration,
		),
		RefreshTokenManager: token.NewManager(
			c.Clock,
			c.SigningKeyAuthentication,
			refreshTokenDuration,
		),
	}
}

type Credentials struct {
	Username string
	Password string
}

type Authentication struct {
	AuthToken    string
	RefreshToken string
}

func (m *Manager) Authenticate(ctx context.Context, c Credentials) (*Authentication, error) {
	u, err := m.userRepo.GetByUsername(ctx, c.Username)
	if err != nil {
		return nil, err
	}

	if !password.Check(c.Password, u.Password) {
		return nil, errors.New("could not authenticate")
	}

	authToken, err := m.AuthTokenManager.New(u.ID.String())
	if err != nil {
		return nil, err
	}

	refreshToken, err := m.RefreshTokenManager.New(u.ID.String())
	if err != nil {
		return nil, err
	}

	return &Authentication{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}, nil
}
