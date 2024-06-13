package app

import (
	"context"
	"log/slog"

	"github.com/extreme-business/lingo/cmd/account/bootstrapping"
	"github.com/extreme-business/lingo/cmd/account/domain"
	"github.com/extreme-business/lingo/cmd/account/user/authentication"
	"github.com/extreme-business/lingo/cmd/account/user/registration"
)

type Account struct {
	logger              *slog.Logger
	bootstrapping       *bootstrapping.Initializer
	accountManager      *authentication.Manager
	registrationManager *registration.Manager
}

func New(
	logger *slog.Logger,
	bootstrapping *bootstrapping.Initializer,
	accountManager *authentication.Manager,
	registrationManager *registration.Manager,
) *Account {
	return &Account{
		logger:              logger,
		bootstrapping:       bootstrapping,
		accountManager:      accountManager,
		registrationManager: registrationManager,
	}
}

func (r *Account) Init(ctx context.Context) error {
	return r.bootstrapping.Setup(ctx)
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
