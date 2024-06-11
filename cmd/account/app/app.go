package app

import (
	"context"
	"log/slog"

	"github.com/dwethmar/lingo/cmd/account/bootstrapping"
	"github.com/dwethmar/lingo/cmd/account/domain"
	"github.com/dwethmar/lingo/cmd/account/user/authentication"
	"github.com/dwethmar/lingo/cmd/account/user/registration"
	"github.com/google/uuid"
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

func (r *Account) CreateUser(ctx context.Context, organizationID uuid.UUID, displayName, email, password string) (*domain.User, error) {
	r.logger.Info("Register")

	user, err := r.registrationManager.Register(ctx, registration.Registration{
		OrganizationID: organizationID,
		DisplayName:    displayName,
		Email:          email,
		Password:       password,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
