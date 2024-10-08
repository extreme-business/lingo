package server

import (
	"context"
	"errors"

	"github.com/extreme-business/lingo/apps/account/app"
	"github.com/extreme-business/lingo/apps/account/domain"
	"github.com/extreme-business/lingo/pkg/grpcerrors"
	"github.com/extreme-business/lingo/pkg/resource"
	"github.com/extreme-business/lingo/pkg/validate"
	protoaccount "github.com/extreme-business/lingo/proto/gen/go/public/account/v1"
	"github.com/google/uuid"
)

type Server struct {
	protoaccount.UnimplementedAccountServiceServer
	account        *app.App
	resourceParser *resource.Parser
}

func New(account *app.App, resourceParser *resource.Parser) *Server {
	return &Server{
		account:        account,
		resourceParser: resourceParser,
	}
}

func (s *Server) CreateUser(ctx context.Context, req *protoaccount.CreateUserRequest) (*protoaccount.CreateUserResponse, error) {
	if req.GetParent() == "" {
		return nil, grpcerrors.NewFieldViolationErr("parent", []grpcerrors.FieldViolation{
			{
				Field:       "parent",
				Description: "parent is required",
			},
		})
	}

	parent, err := s.resourceParser.Parse(req.GetParent())
	if err != nil {
		return nil, err
	}

	var orgID uuid.UUID
	if org := parent.Find(domain.OrganizationCollection); org != nil {
		orgID, err = org.UUID()
		if err != nil {
			return nil, err
		}
	}

	userIn := req.GetUser()
	if userIn == nil {
		return nil, grpcerrors.NewFieldViolationErr("user", []grpcerrors.FieldViolation{
			{
				Field:       "user",
				Description: "user is required",
			},
		})
	}

	user, err := s.account.RegisterUser(ctx, app.RegisterUser{
		OrganizationID: orgID,
		DisplayName:    userIn.GetDisplayName(),
		Email:          userIn.GetEmail(),
		Password:       userIn.GetPassword(),
	})
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

	var userOut protoaccount.User
	if err = user.ToProto(&userOut); err != nil {
		return nil, err
	}

	return &protoaccount.CreateUserResponse{
		User: &userOut,
	}, nil
}

func (s *Server) LoginUser(ctx context.Context, req *protoaccount.LoginUserRequest) (*protoaccount.LoginUserResponse, error) {
	login, err := s.account.LoginUser(
		ctx,
		req.GetEmail(),
		req.GetPassword(),
	)
	if err != nil {
		if errors.Is(err, app.ErrUserNotFound) {
			return nil, grpcerrors.NewNotFoundErr("user not found")
		}
		return nil, err
	}

	var user protoaccount.User
	if err = login.User.ToProto(&user); err != nil {
		return nil, err
	}

	return &protoaccount.LoginUserResponse{
		User:         &user,
		AccessToken:  login.AccessToken,
		RefreshToken: login.RefreshToken,
	}, nil
}

func (s *Server) RefreshToken(ctx context.Context, req *protoaccount.RefreshTokenRequest) (*protoaccount.RefreshTokenResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *Server) ListUsers(ctx context.Context, req *protoaccount.ListUsersRequest) (*protoaccount.ListUsersResponse, error) {
	users, err := s.account.ListUsers(ctx, 0)
	if err != nil {
		return nil, err
	}
	var usersOut []*protoaccount.User
	for _, user := range users {
		var userOut protoaccount.User
		if err = user.ToProto(&userOut); err != nil {
			return nil, err
		}
		usersOut = append(usersOut, &userOut)
	}
	return &protoaccount.ListUsersResponse{
		Users: usersOut,
	}, nil
}
