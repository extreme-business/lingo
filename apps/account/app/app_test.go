// Package app represents the a set of functionality that the account domain provides.
package app

import (
	"context"
	"log/slog"
	"reflect"
	"testing"

	"github.com/extreme-business/lingo/apps/account/auth/authentication"
	"github.com/extreme-business/lingo/apps/account/auth/registration"
	"github.com/extreme-business/lingo/apps/account/domain"
	"github.com/extreme-business/lingo/apps/account/domain/user"
)

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		c := Config{
			Logger:              slog.Default(),
			UserReader:          user.NewReader(nil),
			Authenticator:       authentication.NewManager(authentication.Config{}),
			RegistrationManager: registration.NewManager(registration.Config{}),
		}

		got, err := New(c)
		if err != nil {
			t.Errorf("New() error = %v", err)
			return
		}

		if got == nil {
			t.Errorf("New() = nil")
			return
		}
	})
}

func TestApp_RegisterUser(t *testing.T) {
	type fields struct {
		logger              *slog.Logger
		userReader          *user.Reader
		authenticator       *authentication.Authenticator
		registrationManager *registration.Manager
	}
	type args struct {
		ctx context.Context
		i   RegisterUser
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &App{
				logger:              tt.fields.logger,
				userReader:          tt.fields.userReader,
				authenticator:       tt.fields.authenticator,
				registrationManager: tt.fields.registrationManager,
			}
			got, err := r.RegisterUser(tt.args.ctx, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("App.RegisterUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_ListUsers(t *testing.T) {
	type fields struct {
		logger              *slog.Logger
		userReader          *user.Reader
		authenticator       *authentication.Authenticator
		registrationManager *registration.Manager
	}
	type args struct {
		ctx  context.Context
		page uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &App{
				logger:              tt.fields.logger,
				userReader:          tt.fields.userReader,
				authenticator:       tt.fields.authenticator,
				registrationManager: tt.fields.registrationManager,
			}
			got, err := r.ListUsers(tt.args.ctx, tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.ListUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("App.ListUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}
