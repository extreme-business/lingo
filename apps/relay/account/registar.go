package account

import "github.com/dwethmar/lingo/apps/relay/token"

type Service struct {
	TokenManager *token.Manager
}

func (t *Service) Create(email string) error {
	return t.TokenManager.Create(email)
}

func (t *Service) Register(token string) error {
	return nil
}

func New(tokenManager *token.Manager) *Service {
	return &Service{
		TokenManager: tokenManager,
	}
}
