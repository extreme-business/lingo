package server

import (
	"context"
	"errors"

	"github.com/dwethmar/lingo/cmd/auth/app"
	"github.com/dwethmar/lingo/pkg/grpcerrors"
	"github.com/dwethmar/lingo/pkg/validate"
	protoauth "github.com/dwethmar/lingo/proto/gen/go/public/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	user, err := s.auth.CreateUser(
		ctx,
		req.GetUsername(),
		req.GetEmail(),
		req.GetPassword(),
	)
	if err != nil {
		var vErr *validate.Error
		if errors.As(err, &vErr) {
			return nil, grpcerrors.NewFieldViolationErr("validation error", []grpcerrors.FieldViolation{
				{
					Field:       vErr.Field(),
					Description: vErr.Error(),
				},
			})
		}

		return nil, err
	}

	var userout protoauth.User
	if err = user.ToProto(&userout); err != nil {
		return nil, err
	}

	return &protoauth.CreateUserResponse{
		User: &userout,
	}, nil
}

func (s *Service) LoginUser(ctx context.Context, req *protoauth.LoginUserRequest) (*protoauth.LoginUserResponse, error) {
	login, err := s.auth.LoginUser(
		ctx,
		req.GetEmail(),
		req.GetPassword(),
	)
	if err != nil {
		return nil, err
	}

	var user protoauth.User
	login.User.ToProto(&user)

	grpc.SendHeader(ctx, metadata.Pairs(
		"token", login.Token,
		"refresh_token", login.RefreshToken,
	))

	return &protoauth.LoginUserResponse{
		User: &user,
	}, nil
}
