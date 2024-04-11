package server

import (
	"context"

	"github.com/dwethmar/lingo/apps/auth"

	protoauth "github.com/dwethmar/lingo/protogen/go/proto/public/auth/v1"
)

type Service struct {
	protoauth.UnimplementedAuthServiceServer
	auth *auth.Auth
}

func New(auth *auth.Auth) *Service {
	return &Service{
		auth: auth,
	}
}

func (s *Service) Register(ctx context.Context, req *protoauth.RegisterRequest) (*protoauth.RegisterResponse, error) {
	_, err := s.auth.Register(ctx, "", req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &protoauth.RegisterResponse{}, nil
}
