package app

import (
	"context"
	"log/slog"

	"github.com/extreme-business/lingo/apps/account/auth/authentication"
	"github.com/extreme-business/lingo/apps/account/auth/registration"
	"github.com/extreme-business/lingo/apps/account/domain"
)

type Account struct {
	logger              *slog.Logger
	accountManager      *authentication.Manager
	registrationManager *registration.Manager
}

func New(
	logger *slog.Logger,
	accountManager *authentication.Manager,
	registrationManager *registration.Manager,
) *Account {
	return &Account{
		logger:              logger,
		accountManager:      accountManager,
		registrationManager: registrationManager,
	}
}

func (r *Account) CreateUser(ctx context.Context, u *domain.User, password string) (*domain.User, error) {
	r.logger.Info("Register")

	user, err := r.registrationManager.Register(ctx, registration.Registration{
		User:     u,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
