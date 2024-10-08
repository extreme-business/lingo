// Package app represents the a set of functionality for serving a CMS.
package app

import (
	"context"
	"fmt"

	accountproto "github.com/extreme-business/lingo/proto/gen/go/public/account/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// App is the account manager.
type App struct {
	client accountproto.AccountServiceClient
}

func New(client accountproto.AccountServiceClient) *App {
	return &App{
		client: client,
	}
}

// SuccessResponse is the success response struct.
type SuccessResponse struct {
	AccessToken  string
	RefreshToken string
}

func (a *App) AuthenticateUser(ctx context.Context, email, password string) (*SuccessResponse, error) {
	r, err := a.client.LoginUser(ctx, &accountproto.LoginUserRequest{
		Login: &accountproto.LoginUserRequest_Email{
			Email: email,
		},
		Password: password,
	})
	if err != nil {
		if status.Code(err) == codes.Unauthenticated {
			return nil, fmt.Errorf("invalid email or password")
		}

		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	return &SuccessResponse{
		AccessToken:  r.GetAccessToken(),
		RefreshToken: r.GetRefreshToken(),
	}, nil
}

// Registration is the registration struct.
type Registration struct {
	OrganizationID string
	Email          string
	Password       string
}

func (a *App) RegisterUser(ctx context.Context, r Registration) error {
	_, err := a.client.CreateUser(ctx, &accountproto.CreateUserRequest{
		Parent: fmt.Sprintf("organizations/%s", r.OrganizationID),
		User: &accountproto.User{
			Email:    r.Email,
			Password: r.Password,
		},
	})

	return err
}

func (a *App) ListUsers(ctx context.Context) ([]User, error) {
	r, err := a.client.ListUsers(ctx, &accountproto.ListUsersRequest{})
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(r.GetUsers()))
	for _, u := range r.GetUsers() {
		users = append(users, User{
			Email: u.GetEmail(),
		})
	}

	return users, nil
}
