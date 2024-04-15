package server

import (
	"context"

	"github.com/dwethmar/lingo/cmd/auth/app"
	protoauth "github.com/dwethmar/lingo/proto/gen/go/public/auth/v1"
)

type Service struct {
	protoauth.UnimplementedAuthServiceServer
	auth *app.Auth
}

func New(auth *app.Auth) *Service {
	return &Service{
		auth: auth,
	}
}

func (s *Service) CreateUser(ctx context.Context, req *protoauth.CreateUserRequest) (*protoauth.CreateUserResponse, error) {
	user, err := s.auth.CreateUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &protoauth.CreateUserResponse{
		User: &protoauth.User{
			Id:    user.ID.String(),
			Email: user.Email,
		},
	}, nil
}
