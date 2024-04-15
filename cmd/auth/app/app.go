package app

import (
	"context"
	"log/slog"

	"github.com/dwethmar/lingo/cmd/auth/app/domain"
	"github.com/dwethmar/lingo/cmd/auth/app/token"
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

func (r *Auth) CreateUser(ctx context.Context, username, email, password string) (*domain.User, error) {
	r.logger.Info("Register")

	return &domain.User{}, nil
}
