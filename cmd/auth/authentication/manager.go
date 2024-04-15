package authentication

import (
	"context"
	"time"

	"github.com/dwethmar/lingo/cmd/auth/domain/user"
	"github.com/dwethmar/lingo/cmd/auth/token"
	"github.com/dwethmar/lingo/pkg/clock"
)

type Manager struct {
	repo                       user.Repository
	registrationTokenManager   *token.Manager
	authenticationTokenManager *token.Manager
}

type Config struct {
	Clock                    *clock.Clock
	SigningKeyRegistration   []byte
	SigningKeyAuthentication []byte
	UserRepo                 user.Repository
}

func NewManager(c Config) *Manager {
	return &Manager{
		repo: c.UserRepo,
		registrationTokenManager: token.NewManager(
			c.Clock,
			c.SigningKeyRegistration,
			15*time.Minute,
		),
		authenticationTokenManager: token.NewManager(
			c.Clock,
			c.SigningKeyAuthentication,
			5*time.Minute,
		),
	}
}

func (m *Manager) Login(ctx context.Context, username, password string) (*user.User, error) {
	return &user.User{}, nil
}
