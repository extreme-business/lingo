package app_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/extreme-business/lingo/apps/account/app"
	"github.com/extreme-business/lingo/apps/account/auth/authentication"
	"github.com/extreme-business/lingo/apps/account/auth/registration"
	"github.com/extreme-business/lingo/apps/account/domain"
	"github.com/extreme-business/lingo/apps/account/domain/user"
	"github.com/extreme-business/lingo/apps/account/storage"
	"github.com/extreme-business/lingo/apps/account/storage/postgres"
	"github.com/extreme-business/lingo/apps/account/storage/postgres/seed"
	"github.com/extreme-business/lingo/pkg/database/dbtest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func setupTestDB(ctx context.Context, t *testing.T, name string) *dbtest.PostgresContainer {
	t.Helper()
	dbc := dbtest.SetupPostgres(ctx, t, dbtest.SanitizeDBName(name))
	if err := seed.RunMigrations(ctx, t, dbc.ConnectionString); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return dbc
}

func TestApp_LoginUser(t *testing.T) {
	dbc := setupTestDB(context.Background(), t, t.Name())
	db := dbtest.Connect(context.Background(), t, dbc.ConnectionString)
	dbManager := postgres.NewManager(db)

	// seed the database
	seed.Run(t, dbc.ConnectionString, seed.State{
		Organizations: []*storage.Organization{
			{
				ID:         uuid.MustParse("0463e149-e143-4033-b617-7867824deb0d"),
				LegalName:  "bigorg",
				Slug:       "bigorg",
				CreateTime: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		Users: []*storage.User{
			{
				ID:             uuid.MustParse("d58c4b17-9a1c-4853-9bbf-9467df86307e"),
				OrganizationID: uuid.MustParse("0463e149-e143-4033-b617-7867824deb0d"),
				DisplayName:    "biguser",
				Email:          "buguser@test.com",
				HashedPassword: "$2a$12$8QwkQCXa5omOq4KgNC5Wquv8eGikumiWyUdM0SShfQ/oSt6On0AUu", // thisismypassword
				Status:         "active",
				CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	})

	t.Run("should return the user if the credentials are valid", func(t *testing.T) {
		now := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		app, err := app.New(app.Config{
			Logger: slog.Default(),
			Authenticator: authentication.New(authentication.Config{
				Clock:                  func() time.Time { return now },
				SigningKeyAccessToken:  []byte("access"),
				SigningKeyRefreshToken: []byte("refresh"),
				UserReader:             user.NewReader(dbManager.Op().User),
			}),
			UserReader:          user.NewReader(dbManager.Op().User),
			RegistrationManager: registration.NewManager(registration.Config{}),
		})
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		result, err := app.LoginUser(context.Background(), "buguser@test.com", "thisismypassword")
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}

		if result == nil {
			t.Error("Expected a result, but got nil")
		}

		expect := &domain.User{
			ID:             uuid.MustParse("d58c4b17-9a1c-4853-9bbf-9467df86307e"),
			OrganizationID: uuid.MustParse("0463e149-e143-4033-b617-7867824deb0d"),
			DisplayName:    "biguser",
			Email:          "buguser@test.com",
			Status:         "active",
			CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		if diff := cmp.Diff(expect, result.User); diff != "" {
			t.Errorf("Expected user to match, but got diff: %s", diff)
		}

		if result.AccessToken == "" {
			t.Error("Expected an access token, but got an empty string")
		}

		if result.RefreshToken == "" {
			t.Error("Expected a refresh token, but got an empty string")
		}
	})

	t.Run("should return an error if the credentials are invalid", func(t *testing.T) {
		now := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		a, err := app.New(app.Config{
			Logger: slog.Default(),
			Authenticator: authentication.New(authentication.Config{
				Clock:                  func() time.Time { return now },
				SigningKeyAccessToken:  []byte("access"),
				SigningKeyRefreshToken: []byte("refresh"),
				UserReader:             user.NewReader(dbManager.Op().User),
			}),
			UserReader:          user.NewReader(dbManager.Op().User),
			RegistrationManager: registration.NewManager(registration.Config{}),
		})
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if _, err = a.LoginUser(context.Background(), "buguser@test.com", "invalid"); err == nil {
			t.Error("Expected an error, but got nil")
		} else if !errors.Is(err, app.ErrInvalidCredentials) {
			t.Errorf("Expected an ErrInvalidCredentials error, but got %v", err)
		}
	})

	t.Run("should return an error if the user does not exist", func(t *testing.T) {
		now := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

		a, err := app.New(app.Config{
			Logger: slog.Default(),
			Authenticator: authentication.New(authentication.Config{
				Clock:                  func() time.Time { return now },
				SigningKeyAccessToken:  []byte("access"),
				SigningKeyRefreshToken: []byte("refresh"),
				UserReader:             user.NewReader(dbManager.Op().User),
			}),
			UserReader:          user.NewReader(dbManager.Op().User),
			RegistrationManager: registration.NewManager(registration.Config{}),
		})
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}

		if _, err = a.LoginUser(context.Background(), "invalid", "invalid"); err == nil {
			t.Error("Expected an error, but got nil")
			return
		} else if !errors.Is(err, app.ErrUserNotFound) {
			t.Errorf("Expected an ErrUserNotFound error, but got %v", err)
		}
	})
}
