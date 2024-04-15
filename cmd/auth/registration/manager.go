package registration

import (
	"context"

	"github.com/dwethmar/lingo/cmd/auth/domain/user"
	"github.com/dwethmar/lingo/pkg/clock"
	"github.com/google/uuid"
)

type Manager struct {
	clock *clock.Clock
	repo  user.Repository
}

type Config struct {
	Clock    *clock.Clock
	UserRepo user.Repository
}

func NewManager(c Config) *Manager {
	return &Manager{
		clock: c.Clock,
		repo:  c.UserRepo,
	}
}

// CreateUser creates a new user
func (m *Manager) Register(ctx context.Context, username, email, password string) (*user.User, error) {
	return m.repo.Create(ctx, &user.User{
		ID:         user.ID(uuid.New()),
		Username:   username,
		Email:      email,
		Password:   password,
		CreateTime: m.clock.Now(),
	})
}
