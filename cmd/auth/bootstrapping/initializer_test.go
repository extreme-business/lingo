package bootstrapping_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/dwethmar/lingo/cmd/auth/bootstrapping"
	"github.com/dwethmar/lingo/cmd/auth/storage"
	"github.com/dwethmar/lingo/cmd/auth/storage/postgres"
	"github.com/dwethmar/lingo/cmd/auth/storage/postgres/seed"
	"github.com/dwethmar/lingo/pkg/clock"
	"github.com/dwethmar/lingo/pkg/database"
	"github.com/dwethmar/lingo/pkg/database/dbtest"
	dbmock "github.com/dwethmar/lingo/pkg/database/mock"
	"github.com/dwethmar/lingo/pkg/validate"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		if bootstrapping.New(bootstrapping.Config{}) == nil {
			t.Errorf("New() = nil")
		}
	})
}

func TestInitializer_SetupValidation(t *testing.T) {
	db := database.NewDBWithHandler(&dbmock.DBHandler{
		BeginTxFunc: func(_ context.Context, _ *sql.TxOptions) (*database.Tx, error) {
			return database.NewTxWithHandler(&dbmock.TxHandler{
				RollbackFunc: func() error { return nil },
			}), nil
		},
	})

	type args struct {
		SystemUserConfig         bootstrapping.SystemUserConfig
		SystemOrganizationConfig bootstrapping.SystemOrgConfig
		Clock                    clock.Now
	}
	tests := []struct {
		name      string
		args      args
		want      error
		wantField string
	}{
		{
			name: "should return an error if the user id is zero",
			args: args{
				SystemUserConfig: bootstrapping.SystemUserConfig{
					ID:    uuid.Nil,
					Email: "test@test.com",
				},
				SystemOrganizationConfig: bootstrapping.SystemOrgConfig{
					ID:        uuid.MustParse("2fc54df3-30d3-4b21-8dfd-f4076fc1da65"),
					LegalName: "Test Organization",
				},
				Clock: clock.New(time.UTC, func() time.Time { return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC) }),
			},
			want:      validate.ErrUUIDIsNil,
			wantField: "ID",
		},
		{
			name: "should return an error if the user email is empty",
			args: args{
				SystemUserConfig: bootstrapping.SystemUserConfig{
					ID:    uuid.MustParse("7fb3d880-1db0-464e-b062-a9896cb9bf6c"),
					Email: "",
				},
				SystemOrganizationConfig: bootstrapping.SystemOrgConfig{
					ID:        uuid.MustParse("2fc54df3-30d3-4b21-8dfd-f4076fc1da65"),
					LegalName: "Test Organization",
				},
				Clock: clock.New(time.UTC, func() time.Time { return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC) }),
			},
			want:      validate.ErrStringMinLength,
			wantField: "Email",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			manager := database.NewManager(db, func(_ database.Conn) storage.Repositories { return storage.Repositories{} })

			initializer := bootstrapping.New(bootstrapping.Config{
				SystemUserConfig:         tt.args.SystemUserConfig,
				SystemOrganizationConfig: tt.args.SystemOrganizationConfig,
				Clock:                    tt.args.Clock,
				DBManager:                manager,
			})

			err := initializer.Setup(ctx)

			vErr := &validate.Error{}
			if errors.As(err, &vErr) {
				if vErr.Field() != tt.wantField {
					t.Errorf("expected error to be on field %s, got %s", tt.wantField, vErr.Field())
				}
			} else {
				t.Errorf("expected error to be a validate.Error, got %v", err)
			}
		})
	}
}

func TestInitializer_Setup(t *testing.T) {
	dbc := dbtest.SetupPostgres(t, "bootstrapping_test", seed.RunMigrations)

	seed.Run(t, dbc.ConnectionString, seed.State{
		Organizations: []*storage.Organization{
			seed.NewOrganization(
				"c105ca54-68f0-4bc4-aca1-b54065b4e9b4",
				"Test Organization",
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			),
		},
		Users: []*storage.User{
			seed.NewUser(
				"7fb3d880-1db0-464e-b062-a9896cb9bf6c",
				"c105ca54-68f0-4bc4-aca1-b54065b4e9b4",
				"system",
				"system@system.com",
				"password123",
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			),
		},
	})

	t.Run("Setup", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)
		dbManager := postgres.NewManager(database.NewDB(db))

		initializer := bootstrapping.New(bootstrapping.Config{
			SystemUserConfig: bootstrapping.SystemUserConfig{
				ID:       uuid.MustParse("7fb3d880-1db0-464e-b062-a9896cb9bf6c"),
				Email:    "test@test.com",
				Password: "password",
			},
			SystemOrganizationConfig: bootstrapping.SystemOrgConfig{
				ID:        uuid.MustParse("c105ca54-68f0-4bc4-aca1-b54065b4e9b4"),
				LegalName: "Test Organization",
			},
			Clock:     clock.New(time.UTC, func() time.Time { return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC) }),
			DBManager: dbManager,
		})

		if err := initializer.Setup(ctx); err != nil {
			t.Errorf("Setup() error = %v, want %v", err, nil)
		}

		// Check if the system user and organization were created
		repos := dbManager.Op()

		o, err := repos.Organization.Get(ctx, uuid.MustParse("c105ca54-68f0-4bc4-aca1-b54065b4e9b4"))
		if err != nil {
			t.Errorf("Organization.Get() error = %v, want %v", err, nil)
		}

		expectedOrganization := seed.NewOrganization(
			"c105ca54-68f0-4bc4-aca1-b54065b4e9b4",
			"Test Organization",
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		)

		if diff := cmp.Diff(expectedOrganization, o); diff != "" {
			t.Errorf("Organization.Get() mismatch (-want +got):\n%s", diff)
		}

		u, err := repos.User.Get(ctx, uuid.MustParse("7fb3d880-1db0-464e-b062-a9896cb9bf6c"))
		if err != nil {
			t.Errorf("User.Get() error = %v, want %v", err, nil)
		}

		expectedUser := seed.NewUser(
			"7fb3d880-1db0-464e-b062-a9896cb9bf6c",
			"c105ca54-68f0-4bc4-aca1-b54065b4e9b4",
			"system",
			"test@test.com",
			"",
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		)

		if diff := cmp.Diff(expectedUser, u); diff != "" {
			t.Errorf("User.Get() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Setup with existing organization", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)
		dbManager := postgres.NewManager(database.NewDB(db))

		initializer := bootstrapping.New(bootstrapping.Config{
			SystemUserConfig: bootstrapping.SystemUserConfig{
				ID:       uuid.MustParse("7fb3d880-1db0-464e-b062-a9896cb9bf6c"),
				Email:    "test@test.nl",
				Password: "password",
			},
			SystemOrganizationConfig: bootstrapping.SystemOrgConfig{
				ID:        uuid.MustParse("c105ca54-68f0-4bc4-aca1-b54065b4e9b4"),
				LegalName: "Test Organization",
			},
			Clock:     clock.New(time.UTC, func() time.Time { return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC) }),
			DBManager: dbManager,
		})

		if err := initializer.Setup(ctx); err != nil {
			t.Errorf("Setup() error = %v, want %v", err, nil)
		}
	})
}
