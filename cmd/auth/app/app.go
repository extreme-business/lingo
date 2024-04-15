package app

import (
	"context"
	"log/slog"

	"github.com/dwethmar/lingo/apps/relay/token"
	"github.com/dwethmar/lingo/cmd/auth/app/domain"
)

type Auth struct {
	logger *slog.Logger
}

func New(
	logger *slog.Logger,
	registrationTokenManager *token.Manager,
	authenticationTokenManager *token.Manager,
) *Auth {
	return &Auth{
		logger: logger,
	}
}

func (r *Auth) Register(ctx context.Context, username, email, password string) (*domain.User, error) {
	r.logger.Info("Register")

	return &domain.User{}, nil
}

func (r *Auth) CreateMessage(ctx context.Context, message string) error {
	r.logger.Info("CreateMessage")

	return nil
}
