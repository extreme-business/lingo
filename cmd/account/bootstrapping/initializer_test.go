package bootstrapping_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/extreme-business/lingo/cmd/account/bootstrapping"
	"github.com/extreme-business/lingo/cmd/account/storage"
	"github.com/extreme-business/lingo/cmd/account/storage/postgres"
	"github.com/extreme-business/lingo/cmd/account/storage/postgres/seed"
	"github.com/extreme-business/lingo/pkg/clock"
	"github.com/extreme-business/lingo/pkg/database"
	"github.com/extreme-business/lingo/pkg/database/dbtest"
	dbmock "github.com/extreme-business/lingo/pkg/database/mock"
	"github.com/extreme-business/lingo/pkg/validate"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		b, err := bootstrapping.New(bootstrapping.Config{
			Logger: slog.Default(),
			SystemUserConfig: bootstrapping.SystemUserConfig{
				ID:       uuid.MustParse("7fb3d880-1db0-464e-b062-a9896cb9bf6c"),
				Email:    "test@test.nl",
				Password: "password",
			},
			SystemOrganizationConfig: bootstrapping.SystemOrgConfig{
				ID:        uuid.MustParse("c105ca54-68f0-4bc4-aca1-b54065b4e9b4"),
				LegalName: "Test Organization",
				Slug:      "test-organization",
			},
			Clock: clock.Default(),
			DBManager: database.NewManager(database.NewDBWithHandler(&dbmock.DBHandler{}), func(_ database.Conn) storage.Repositories {
				return storage.Repositories{}
			}),
		})
		if err != nil {
			t.Errorf("New() error = %v, want %v", err, nil)
		}

		if b == nil {
			t.Errorf("New() got nil, want not nil")
		}
	})
}

func TestInitializer_NewManager(t *testing.T) {
	t.Run("NewManager", func(t *testing.T) {
		b, err := bootstrapping.New(bootstrapping.Config{
			Logger: slog.Default(),
			SystemUserConfig: bootstrapping.SystemUserConfig{
				ID:       uuid.MustParse("7fb3d880-1db0-464e-b062-a9896cb9bf6c"),
				Email:    "test@test.nl",
				Password: "password",
			},
			SystemOrganizationConfig: bootstrapping.SystemOrgConfig{
				ID:        uuid.MustParse("c105ca54-68f0-4bc4-aca1-b54065b4e9b4"),
				LegalName: "Test Organization",
				Slug:      "test-organization",
			},
			Clock: clock.Default(),
			DBManager: database.NewManager(database.NewDBWithHandler(&dbmock.DBHandler{}), func(_ database.Conn) storage.Repositories {
				return storage.Repositories{}
			}),
		})

		if err != nil {
			t.Errorf("New() error = %v, want %v", err, nil)
		}

		if b == nil {
			t.Errorf("New() got nil, want not nil")
		}

	})

}

func TestInitializer_Setup(t *testing.T) {
	dbc := dbtest.SetupPostgres(t, "bootstrapping_test", seed.RunMigrations)

	seed.Run(t, dbc.ConnectionString, seed.State{
		Organizations: []*storage.Organization{
			seed.NewOrganization(
				"c105ca54-68f0-4bc4-aca1-b54065b4e9b4",
				"Test Organization",
				"test-organization",
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

		initializer, err := bootstrapping.New(bootstrapping.Config{
			Logger: slog.Default(),
			SystemUserConfig: bootstrapping.SystemUserConfig{
				ID:       uuid.MustParse("7fb3d880-1db0-464e-b062-a9896cb9bf6c"),
				Email:    "test@test.com",
				Password: "password",
			},
			SystemOrganizationConfig: bootstrapping.SystemOrgConfig{
				ID:        uuid.MustParse("c105ca54-68f0-4bc4-aca1-b54065b4e9b4"),
				LegalName: "Test Organization",
				Slug:      "test-organization",
			},
			Clock:     clock.New(time.UTC, func() time.Time { return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC) }),
			DBManager: dbManager,
		})

		if err != nil {
			t.Errorf("New() error = %v, want %v", err, nil)
		}

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
			"test-organization",
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

		initializer, err := bootstrapping.New(bootstrapping.Config{
			Logger: slog.Default(),
			SystemUserConfig: bootstrapping.SystemUserConfig{
				ID:       uuid.MustParse("7fb3d880-1db0-464e-b062-a9896cb9bf6c"),
				Email:    "test@test.nl",
				Password: "password",
			},
			SystemOrganizationConfig: bootstrapping.SystemOrgConfig{
				ID:        uuid.MustParse("c105ca54-68f0-4bc4-aca1-b54065b4e9b4"),
				LegalName: "Test Organization",
				Slug:      "test-organization",
			},
			Clock:     clock.New(time.UTC, func() time.Time { return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC) }),
			DBManager: dbManager,
		})

		if err != nil {
			t.Errorf("New() error = %v, want %v", err, nil)
		}

		if err := initializer.Setup(ctx); err != nil {
			t.Errorf("Setup() error = %v, want %v", err, nil)
		}
	})
}

func TestSystemUserConfig_Validate(t *testing.T) {
	type fields struct {
		ID       uuid.UUID
		Email    string
		Password string
	}
	tests := []struct {
		name     string
		fields   fields
		err      error
		errField string
	}{
		{
			name: "should return no error if the config is valid",
			fields: fields{
				ID:       uuid.MustParse("7fb3d880-1db0-464e-b062-a9896cb9bf6c"),
				Email:    "test@test.nl",
				Password: "Password123@",
			},
			err:      nil,
			errField: "",
		},
		{
			name: "should return an error if the ID is empty",
			fields: fields{
				ID:       uuid.Nil,
				Email:    "test@test.nl",
				Password: "Password123@",
			},
			err:      validate.ErrUUIDIsNil,
			errField: "ID",
		},
		{
			name: "should return an error if the email is empty",
			fields: fields{
				ID:       uuid.MustParse("7fb3d880-1db0-464e-b062-a9896cb9bf6c"),
				Email:    "",
				Password: "Password123@",
			},
			err:      validate.ErrStringMinLength,
			errField: "Email",
		},
		{
			name: "should return an error if the password is empty",
			fields: fields{
				ID:       uuid.MustParse("7fb3d880-1db0-464e-b062-a9896cb9bf6c"),
				Email:    "test@test.nl",
				Password: "",
			},
			err:      validate.ErrStringMinLength,
			errField: "Password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &bootstrapping.SystemUserConfig{
				ID:       tt.fields.ID,
				Email:    tt.fields.Email,
				Password: tt.fields.Password,
			}

			err := c.Validate()
			if !errors.Is(err, tt.err) {
				t.Errorf("SystemUserConfig.Validate() error = %v, wantErr %v", err, tt.err)
			}

			if v, ok := validate.ToError(err); ok {
				if v.Field() != tt.errField {
					t.Errorf("SystemUserConfig.Validate() error field = %v, want %v", v.Field(), tt.errField)
				}
			}
		})
	}
}

func TestSystemOrgConfig_Validate(t *testing.T) {
	type fields struct {
		ID        uuid.UUID
		LegalName string
		Slug      string
	}
	tests := []struct {
		name     string
		fields   fields
		err      error
		errField string
	}{
		{
			name: "should return no error if the config is valid",
			fields: fields{
				ID:        uuid.MustParse("c105ca54-68f0-4bc4-aca1-b54065b4e9b4"),
				LegalName: "Test Organization",
				Slug:      "test-organization",
			},
			err:      nil,
			errField: "",
		},
		{
			name: "should return an error if the ID is empty",
			fields: fields{
				ID:        uuid.Nil,
				LegalName: "Test Organization",
				Slug:      "test-organization",
			},
			err:      validate.ErrUUIDIsNil,
			errField: "ID",
		},
		{
			name: "should return an error if the legal name is empty",
			fields: fields{
				ID:        uuid.MustParse("c105ca54-68f0-4bc4-aca1-b54065b4e9b4"),
				LegalName: "",
				Slug:      "test-organization",
			},
			err:      validate.ErrStringRequired,
			errField: "LegalName",
		},
		{
			name: "should return an error if the slug is empty",
			fields: fields{
				ID:        uuid.MustParse("c105ca54-68f0-4bc4-aca1-b54065b4e9b4"),
				LegalName: "Test Organization",
				Slug:      "",
			},
			err:      validate.ErrStringRequired,
			errField: "Slug",
		},
		{
			name: "should return an error if the slug is too long",
			fields: fields{
				ID:        uuid.MustParse("c105ca54-68f0-4bc4-aca1-b54065b4e9b4"),
				LegalName: "Test Organization",
				Slug:      "test-test-test-test-test-test-test-test-test-test-test-test-test-test-test-test-test-test-test-test-test-test-test-test",
			},
			err:      validate.ErrStringMaxLength,
			errField: "Slug",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &bootstrapping.SystemOrgConfig{
				ID:        tt.fields.ID,
				LegalName: tt.fields.LegalName,
				Slug:      tt.fields.Slug,
			}

			err := c.Validate()
			if !errors.Is(err, tt.err) {
				t.Errorf("SystemOrgConfig.Validate() error = %v, wantErr %v", err, tt.err)
			}

			if v, ok := validate.ToError(err); ok {
				if v.Field() != tt.errField {
					t.Errorf("SystemOrgConfig.Validate() error field = %v, want %v", v.Field(), tt.errField)
				}
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	type fields struct {
		Logger                   *slog.Logger
		SystemUserConfig         bootstrapping.SystemUserConfig
		SystemOrganizationConfig bootstrapping.SystemOrgConfig
		Clock                    clock.Now
		DBManager                storage.DBManager
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := bootstrapping.Config{
				Logger:                   tt.fields.Logger,
				SystemUserConfig:         tt.fields.SystemUserConfig,
				SystemOrganizationConfig: tt.fields.SystemOrganizationConfig,
				Clock:                    tt.fields.Clock,
				DBManager:                tt.fields.DBManager,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitializer_setup(t *testing.T) {

}
