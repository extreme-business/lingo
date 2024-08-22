package bootstrapping_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/extreme-business/lingo/apps/account/bootstrapping"
	"github.com/extreme-business/lingo/apps/account/storage"
	"github.com/extreme-business/lingo/apps/account/storage/postgres"
	"github.com/extreme-business/lingo/apps/account/storage/postgres/seed"
	"github.com/extreme-business/lingo/pkg/database"
	"github.com/extreme-business/lingo/pkg/database/dbtest"
	dbmock "github.com/extreme-business/lingo/pkg/database/mock"
	"github.com/extreme-business/lingo/pkg/validate"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func setupTestDB(ctx context.Context, t *testing.T, name string) *dbtest.PostgresContainer {
	t.Helper()
	// replace special characters with underscores
	dbc := dbtest.SetupPostgres(ctx, t, dbtest.SanitizeDBName(name))
	if err := seed.RunMigrations(ctx, t, dbc.ConnectionString); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return dbc
}

func dbUserDiff(t *testing.T, r storage.UserRepository, u *storage.User) string {
	t.Helper()
	user, err := r.Get(context.Background(), u.ID)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	return cmp.Diff(u, user)
}

func dbOrgDiff(t *testing.T, r storage.OrganizationRepository, o *storage.Organization) string {
	t.Helper()
	org, err := r.Get(context.Background(), o.ID)
	if err != nil {
		t.Fatalf("failed to get organization: %v", err)
	}

	return cmp.Diff(o, org)
}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		b, err := bootstrapping.New(bootstrapping.Config{
			Logger: slog.Default(),
			Clock:  time.Now,
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
			Clock:  time.Now,
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

func seedSystem(t *testing.T, connectionString string) {
	seed.Run(t, connectionString, seed.State{
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
				"active",
				"system@system.com",
				"password123",
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Time{},
			),
		},
	})
}

func TestInitializer_Setup(t *testing.T) {
	t.Run("Setup", func(t *testing.T) {
		ctx := context.Background()
		dbc := setupTestDB(ctx, t, t.Name())
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)
		dbManager := postgres.NewManager(database.NewDB(db))
		seedSystem(t, dbc.ConnectionString)

		initializer, err := bootstrapping.New(bootstrapping.Config{
			Logger:    slog.Default(),
			Clock:     func() time.Time { return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC) },
			DBManager: dbManager,
		})

		if err != nil {
			t.Errorf("New() error = %v, want %v", err, nil)
		}

		suc := bootstrapping.SystemUserConfig{
			ID:       uuid.MustParse("44756c0a-28b7-40c7-a066-f8db23d7dbe3"),
			Email:    "test@test.com",
			Password: "password",
		}

		soc := bootstrapping.SystemOrgConfig{
			ID:        uuid.MustParse("c105ca54-68f0-4bc4-aca1-b54065b4e9b4"),
			LegalName: "Test Organization",
			Slug:      "test-organization",
		}

		if err = initializer.Setup(ctx, suc, soc); err != nil {
			t.Errorf("Setup() error = %v, want %v", err, nil)
		}

		// Check if the system user and organization were created
		repos := dbManager.Op()
		if diff := dbOrgDiff(t, repos.Organization, seed.NewOrganization(
			"c105ca54-68f0-4bc4-aca1-b54065b4e9b4",
			"Test Organization",
			"test-organization",
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		)); diff != "" {
			t.Errorf("Organization.Get() mismatch (-want +got):\n%s", diff)
		}

		if diff := dbUserDiff(t, repos.User, seed.NewUser(
			"44756c0a-28b7-40c7-a066-f8db23d7dbe3",
			"c105ca54-68f0-4bc4-aca1-b54065b4e9b4",
			"system",
			"active",
			"test@test.com",
			"",
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
		)); diff != "" {
			t.Errorf("User.Get() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Setup with an already existing system user and organization should update new changes", func(t *testing.T) {
		ctx := context.Background()
		dbc := setupTestDB(ctx, t, t.Name())
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)
		dbManager := postgres.NewManager(database.NewDB(db))

		initializer, err := bootstrapping.New(bootstrapping.Config{
			Logger:    slog.Default(),
			Clock:     func() time.Time { return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC) },
			DBManager: dbManager,
		})

		if err != nil {
			t.Errorf("New() error = %v, want %v", err, nil)
		}

		suc := bootstrapping.SystemUserConfig{
			ID:       uuid.MustParse("44756c0a-28b7-40c7-a066-f8db23d7dbe3"),
			Email:    "test@test.com",
			Password: "password",
		}

		soc := bootstrapping.SystemOrgConfig{
			ID:        uuid.MustParse("c105ca54-68f0-4bc4-aca1-b54065b4e9b4"),
			LegalName: "Test Organization",
			Slug:      "test-organization",
		}

		if err = initializer.Setup(ctx, suc, soc); err != nil {
			t.Errorf("Setup() error = %v, want %v", err, nil)
		}

		initializer2, err := bootstrapping.New(bootstrapping.Config{
			Logger:    slog.Default(),
			Clock:     func() time.Time { return time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC) },
			DBManager: dbManager,
		})
		if err != nil {
			t.Errorf("New() error = %v, want %v", err, nil)
		}

		suc = bootstrapping.SystemUserConfig{
			ID:       uuid.MustParse("44756c0a-28b7-40c7-a066-f8db23d7dbe3"),
			Email:    "updated-test@test.com",
			Password: "password",
		}

		soc = bootstrapping.SystemOrgConfig{
			ID:        uuid.MustParse("c105ca54-68f0-4bc4-aca1-b54065b4e9b4"),
			LegalName: "updated-Test Organization",
			Slug:      "updated-test-organization",
		}

		if err = initializer2.Setup(ctx, suc, soc); err != nil {
			t.Errorf("Setup() error = %v, want %v", err, nil)
		}

		// Check if the system user and organization were created
		repos := dbManager.Op()

		if diff := dbOrgDiff(t, repos.Organization, seed.NewOrganization(
			"c105ca54-68f0-4bc4-aca1-b54065b4e9b4",
			"updated-Test Organization",
			"updated-test-organization",
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
		)); diff != "" {
			t.Errorf("Organization.Get() mismatch (-want +got):\n%s", diff)
		}

		if diff := dbUserDiff(t, repos.User, seed.NewUser(
			"44756c0a-28b7-40c7-a066-f8db23d7dbe3",
			"c105ca54-68f0-4bc4-aca1-b54065b4e9b4",
			"system",
			"active",
			"updated-test@test.com",
			"",
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
			time.Time{},
		)); diff != "" {
			t.Errorf("User.Get() mismatch (-want +got):\n%s", diff)
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

			if v, ok := validate.AssertError(err); ok {
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
			err:      validate.ErrEmptyString,
			errField: "LegalName",
		},
		{
			name: "should return an error if the slug is empty",
			fields: fields{
				ID:        uuid.MustParse("c105ca54-68f0-4bc4-aca1-b54065b4e9b4"),
				LegalName: "Test Organization",
				Slug:      "",
			},
			err:      validate.ErrEmptyString,
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

			if v, ok := validate.AssertError(err); ok {
				if v.Field() != tt.errField {
					t.Errorf("SystemOrgConfig.Validate() error field = %v, want %v", v.Field(), tt.errField)
				}
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {

}

func TestInitializer_setup(_ *testing.T) {

}
