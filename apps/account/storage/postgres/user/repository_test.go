package user_test

import (
	"context"
	_ "embed"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/extreme-business/lingo/apps/account/storage"
	"github.com/extreme-business/lingo/apps/account/storage/postgres/seed"
	"github.com/extreme-business/lingo/apps/account/storage/postgres/user"
	"github.com/extreme-business/lingo/pkg/database"
	"github.com/extreme-business/lingo/pkg/database/dbtest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func setupTestDB(ctx context.Context, t *testing.T, name string) *dbtest.PostgresContainer {
	t.Helper()
	dbc := dbtest.SetupPostgres(ctx, t, name)
	if err := seed.RunMigrations(ctx, t, dbc.ConnectionString); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return dbc
}

func TestNew(t *testing.T) {
	t.Run("should return a new repository", func(t *testing.T) {
		if got := user.New(nil); got == nil {
			t.Error("expected repository")
		}
	})
}

func TestRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	dbc := setupTestDB(context.Background(), t, "user")

	seed.Run(t, dbc.ConnectionString, seed.State{
		Organizations: []*storage.Organization{
			seed.NewOrganization(
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			),
		},
	})

	t.Run("Create should create a new user", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)

		repo := user.New(database.NewDB(db))
		user, err := repo.Create(ctx, seed.NewUser(
			"35297169-89d8-444d-8499-c6341e3a0770",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			"active",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
		))

		if err != nil {
			t.Fatal(err)
		}

		expect := seed.NewUser(
			"35297169-89d8-444d-8499-c6341e3a0770",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			"active",
			"test@test.com",
			"",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
		)

		if diff := cmp.Diff(expect, user); diff != "" {
			t.Errorf("Create() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("should return an error if the user id already exists", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)

		repo := user.New(database.NewDB(db))
		_, err := repo.Create(ctx, seed.NewUser(
			"485819f0-9e48-4d25-b07b-6de8a2076be2",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"username_300",
			"active",
			"test_300@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
		))

		if err != nil {
			t.Fatal(err)
		}

		_, err = repo.Create(ctx, seed.NewUser(
			"485819f0-9e48-4d25-b07b-6de8a2076be2",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"username_301",
			"active",
			"test_3001@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
		))

		if !errors.Is(err, storage.ErrConflictUserID) {
			t.Errorf("expected %q, got %q", storage.ErrConflictUserID, err)
		}
	})

	t.Run("should return an error if the user email already exists", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)

		repo := user.New(database.NewDB(db))
		_, err := repo.Create(ctx, seed.NewUser(
			"2e56b481-05fe-4ce3-b072-a94fbf8aeab3",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"username_400",
			"active",
			"test_400@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
		))

		if err != nil {
			t.Fatal(err)
		}

		_, err = repo.Create(ctx, seed.NewUser(
			"5e6f2f35-1de1-4803-8fdd-9b67706f887e",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"username_401",
			"active",
			"test_400@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
		))

		if !errors.Is(err, storage.ErrConflictUserEmail) {
			t.Errorf("expected %q, got %q", storage.ErrConflictUserEmail, err)
		}
	})
}

func TestRepository_Get(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	dbc := setupTestDB(context.Background(), t, "user")

	seed.Run(t, dbc.ConnectionString, seed.State{
		Organizations: []*storage.Organization{
			seed.NewOrganization(
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			),
		},
		Users: []*storage.User{
			seed.NewUser(
				"35297169-89d8-444d-8499-c6341e3a0770",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"active",
				"test@test.com",
				"password",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Time{},
			),
		},
	})

	t.Run("Get should get a user", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)

		repo := user.New(database.NewDB(db))

		user, err := repo.Get(ctx, uuid.MustParse("35297169-89d8-444d-8499-c6341e3a0770"))
		if err != nil {
			t.Fatal(err)
		}

		expect := seed.NewUser(
			"35297169-89d8-444d-8499-c6341e3a0770",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			"active",
			"test@test.com",
			"",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
		)

		if diff := cmp.Diff(expect, user); diff != "" {
			t.Errorf("Create() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("should return an error if the user does not exist", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)
		repo := user.New(database.NewDB(db))
		u, err := repo.Get(ctx, uuid.MustParse("946adb15-195e-44df-922b-4a45b9505684"))

		if !errors.Is(err, storage.ErrUserNotFound) {
			t.Errorf("expected %q, got %q", storage.ErrUserNotFound, err)
		}

		if u != nil {
			t.Errorf("expected nil, got %v", u)
		}
	})
}

func TestRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	dbc := setupTestDB(context.Background(), t, "user")

	seed.Run(t, dbc.ConnectionString, seed.State{
		Organizations: []*storage.Organization{
			seed.NewOrganization(
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			),
			seed.NewOrganization(
				"f226487d-61ff-4a18-a2d9-ab888b22dbc8",
				"test2",
				"test2",
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			),
		},
		Users: []*storage.User{
			seed.NewUser(
				"35297169-89d8-444d-8499-c6341e3a0770",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"active",
				"test@test.com",
				"password",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Time{},
			),
		},
	})

	t.Run("should update a user", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)
		recorder := dbtest.NewRecorder(database.NewDB(db))
		repo := user.New(recorder)

		user, err := repo.Update(ctx, seed.NewUser(
			"35297169-89d8-444d-8499-c6341e3a0770",
			"f226487d-61ff-4a18-a2d9-ab888b22dbc8",
			"test2",    // updated username
			"inactive", // updated status
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
		), []storage.UserField{
			storage.UserOrganizationID,
			storage.UserDisplayName,
			storage.UserEmail,
			storage.UserStatus,
			storage.UserHashedPassword,
			storage.UserUpdateTime,
		})

		if err != nil {
			t.Fatal(err)
		}

		expect := seed.NewUser(
			"35297169-89d8-444d-8499-c6341e3a0770",
			"f226487d-61ff-4a18-a2d9-ab888b22dbc8",
			"test2",
			"inactive",
			"test@test.com",
			"",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
		)

		if diff := cmp.Diff(expect, user); diff != "" {
			t.Errorf("Create() mismatch (-want +got):\n%s", diff)
		}

		// check if only the username was updated
		query := recorder.RowQueries[0].Query
		query = strings.TrimSpace(strings.ReplaceAll(query, "\n", " "))
		expectedQuery := "UPDATE users SET organization_id = $1, display_name = $2, email = $3, status = $4, hashed_password = $5, update_time = $6 WHERE id = $7 RETURNING id, organization_id,  display_name, email, status, create_time, update_time, delete_time;"

		if query != expectedQuery {
			t.Errorf("expected %q, got %q", expectedQuery, query)
		}
	})

	t.Run("should return an error if the user does not exist", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)

		repo := user.New(database.NewDB(db))
		_, err := repo.Update(ctx, seed.NewUser(
			"f2e8b3cd-07a3-4d7c-9eef-cf02452d8332",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			"active",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
		), []storage.UserField{
			storage.UserDisplayName,
			storage.UserUpdateTime,
		})

		if !errors.Is(err, storage.ErrUserNotFound) {
			t.Errorf("expected %q, got %q", storage.ErrUserNotFound, err)
		}
	})

	t.Run("should return an error if no fields are provided", func(t *testing.T) {
		ctx := context.Background()
		r, err := user.New(nil).Update(
			ctx,
			seed.NewUser(
				"957b12c5-1071-40d9-8bec-6ed195c8cfbf",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"active",
				"test@test.com",
				"password",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Time{},
			),
			[]storage.UserField{}, // no fields
		)

		if err == nil || !errors.Is(err, storage.ErrNoUserFieldsToUpdate) {
			t.Errorf("expected %q, got %q", storage.ErrNoUserFieldsToUpdate, err)
		}

		if r != nil {
			t.Errorf("expected nil, got %v", r)
		}
	})
}

func TestRepository_Update_fields(t *testing.T) {
	dbc := setupTestDB(context.Background(), t, "user")

	seed.Run(t, dbc.ConnectionString, seed.State{
		Organizations: []*storage.Organization{
			seed.NewOrganization(
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			),
		},
		Users: []*storage.User{
			seed.NewUser(
				"957b12c5-1071-40d9-8bec-6ed195c8cfbf",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"active",
				"test@test.com",
				"password",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Time{},
			),
		},
	})

	tests := []struct {
		name   string
		user   *storage.User
		fields []storage.UserField
		err    error
	}{
		{
			name: "should return an error no fields are provided",
			user: seed.NewUser(
				"957b12c5-1071-40d9-8bec-6ed195c8cfbf",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"active",
				"test@test.com",
				"password",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Time{},
			),
			fields: []storage.UserField{},
			err:    storage.ErrNoUserFieldsToUpdate,
		},
		{
			name: "should return an error field is create_time",
			user: seed.NewUser(
				"957b12c5-1071-40d9-8bec-6ed195c8cfbf",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"active",
				"test@test.com",
				"password",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Time{},
			),
			fields: []storage.UserField{storage.UserCreateTime},
			err:    storage.ErrImmutableUserCreateTime,
		},
		{
			name: "should return an error if the field is unknown",
			user: seed.NewUser(
				"957b12c5-1071-40d9-8bec-6ed195c8cfbf",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"active",
				"test@test.com",
				"password",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Time{},
			),
			fields: []storage.UserField{"unknown"},
			err:    storage.ErrUserUnknownField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := user.New(nil).Update(ctx, tt.user, tt.fields)

			if (err != nil) != errors.Is(err, tt.err) {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.err)
				return
			}
		})
	}
}

func TestRepository_GetByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	dbc := setupTestDB(context.Background(), t, "user")

	seed.Run(t, dbc.ConnectionString, seed.State{
		Organizations: []*storage.Organization{
			seed.NewOrganization(
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			),
		},
		Users: []*storage.User{
			seed.NewUser(
				"82651da9-c2ff-4152-8eae-7555d5a42aad",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"active",
				"test@test.com",
				"password",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Time{},
			),
		},
	})

	t.Run("should get a user by email", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)
		repo := user.New(database.NewDB(db))
		user, err := repo.GetByEmail(ctx, "test@test.com")

		if err != nil {
			t.Fatal(err)
		}

		expect := seed.NewUser(
			"82651da9-c2ff-4152-8eae-7555d5a42aad",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			"active",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
		)

		if diff := cmp.Diff(expect, user); diff != "" {
			t.Errorf("Create() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("should return an error if the user does not exist", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)

		repo := user.New(database.NewDB(db))
		u, err := repo.GetByEmail(ctx, "test2@test.com")

		if err == nil || !errors.Is(err, storage.ErrUserNotFound) {
			t.Errorf("expected %q, got %q", storage.ErrUserNotFound, err)
		}

		if u != nil {
			t.Errorf("expected nil, got %v", u)
		}
	})
}

func TestRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	dbc := setupTestDB(context.Background(), t, "user")

	seed.Run(t, dbc.ConnectionString, seed.State{
		Organizations: []*storage.Organization{
			seed.NewOrganization(
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			),
		},
		Users: []*storage.User{
			seed.NewUser(
				"35297169-89d8-444d-8499-c6341e3a0770",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"active",
				"test@test.com",
				"password",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Time{},
			),
		},
	})

	t.Run("Delete should delete a user", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)
		repo := user.New(database.NewDB(db))
		if err := repo.Delete(ctx, uuid.MustParse("35297169-89d8-444d-8499-c6341e3a0770")); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should return an error if the user does not exist", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)

		repo := user.New(database.NewDB(db))
		err := repo.Delete(ctx, uuid.MustParse("82651da9-c2ff-4152-8eae-7555d5a42aad"))

		if !errors.Is(err, storage.ErrUserNotFound) {
			t.Errorf("expected %q, got %q", storage.ErrUserNotFound, err)
		}
	})
}

func setupTestDatabaseForList(t *testing.T) *dbtest.PostgresContainer {
	t.Helper()
	dbc := setupTestDB(context.Background(), t, "user_list")
	seed.Run(t, dbc.ConnectionString, seed.State{
		Organizations: []*storage.Organization{
			seed.NewOrganization(
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			),
		},
		Users: []*storage.User{
			seed.NewUser(
				"7f62b8ca-6c97-4081-adbc-2b4611b41617",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"active",
				"test@test.com",
				"password",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Time{},
			),
			seed.NewUser(
				"92c02cd8-9286-4687-bb4d-60ee95c769ed",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test2",
				"active",
				"test2@test.com",
				"password",
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Time{},
			),
			seed.NewUser(
				"2198aab6-2ece-429a-8d7c-1654ab8b7d8f",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test3",
				"active",
				"test3@test.com",
				"password",
				time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				time.Time{},
			),
		},
	})
	return dbc
}

func connectAndCreateRepo(t *testing.T, db database.Conn) (*user.Repository, *dbtest.Recorder) {
	t.Helper()
	recorder := dbtest.NewRecorder(db)
	repo := user.New(recorder)
	return repo, recorder
}

func assertUsers(t *testing.T, expect, actual []*storage.User) {
	t.Helper()
	if diff := cmp.Diff(expect, actual); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func assertQuery(t *testing.T, recorder *dbtest.Recorder, expectedQuery string) {
	t.Helper()
	if expectedQuery != "" {
		if len(recorder.Queries) != 1 {
			t.Errorf("expected 1 query, got %d", len(recorder.Queries))
			return // no need to check the query
		}

		query := recorder.Queries[0].Query
		query = strings.TrimSpace(strings.ReplaceAll(query, "\n", " "))
		if query != expectedQuery {
			t.Errorf("expected %q, got %q", expectedQuery, query)
		}
	}
}

func assertError(t *testing.T, expected, actual error) {
	t.Helper()
	if expected == nil && actual != nil {
		t.Errorf("unexpected error: %v", actual)
	}
	if expected != nil && actual == nil {
		t.Errorf("expected error: %v, got nil", expected)
	}
	if expected != nil && !errors.Is(actual, expected) {
		t.Errorf("expected error: %v, got: %v", expected, actual)
	}
}

func TestRepository_List(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	dbc := setupTestDatabaseForList(t)

	testCases := []struct {
		name          string
		ctx           context.Context
		db            database.Conn
		pagination    storage.Pagination
		orderBy       storage.UserOrderBy
		conditions    []storage.Condition
		expectedUsers []*storage.User
		expectedQuery string
		expectedError error
	}{
		{
			name:       "should list users",
			ctx:        context.Background(),
			db:         database.NewDB(dbtest.Connect(context.Background(), t, dbc.ConnectionString)),
			pagination: storage.Pagination{Limit: 0, Offset: 0},
			expectedUsers: []*storage.User{
				seed.NewUser(
					"7f62b8ca-6c97-4081-adbc-2b4611b41617",
					"7bb443e5-8974-44c2-8b7c-b95124205264",
					"test",
					"active",
					"test@test.com",
					"",
					time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Time{},
				),
				seed.NewUser(
					"92c02cd8-9286-4687-bb4d-60ee95c769ed",
					"7bb443e5-8974-44c2-8b7c-b95124205264",
					"test2",
					"active",
					"test2@test.com",
					"",
					time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Time{},
				),
				seed.NewUser(
					"2198aab6-2ece-429a-8d7c-1654ab8b7d8f",
					"7bb443e5-8974-44c2-8b7c-b95124205264",
					"test3",
					"active",
					"test3@test.com",
					"",
					time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Time{},
				),
			},
			expectedQuery: "SELECT u.id, u.organization_id, u.display_name, u.email, u.status, u.create_time, u.update_time, u.delete_time FROM users u;",
		},
		{
			name:       "should list users with organization id predicate",
			ctx:        context.Background(),
			db:         database.NewDB(dbtest.Connect(context.Background(), t, dbc.ConnectionString)),
			pagination: storage.Pagination{Limit: 0, Offset: 0},
			conditions: []storage.Condition{
				storage.UserByOrganizationIDCondition{
					OrganizationID: uuid.MustParse("7bb443e5-8974-44c2-8b7c-b95124205264"),
				},
			},
			expectedUsers: []*storage.User{
				seed.NewUser(
					"7f62b8ca-6c97-4081-adbc-2b4611b41617",
					"7bb443e5-8974-44c2-8b7c-b95124205264",
					"test",
					"active",
					"test@test.com",
					"",
					time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Time{},
				),
				seed.NewUser(
					"92c02cd8-9286-4687-bb4d-60ee95c769ed",
					"7bb443e5-8974-44c2-8b7c-b95124205264",
					"test2",
					"active",
					"test2@test.com",
					"",
					time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Time{},
				),
				seed.NewUser(
					"2198aab6-2ece-429a-8d7c-1654ab8b7d8f",
					"7bb443e5-8974-44c2-8b7c-b95124205264",
					"test3",
					"active",
					"test3@test.com",
					"",
					time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Time{},
				),
			},
			expectedQuery: "SELECT u.id, u.organization_id, u.display_name, u.email, u.status, u.create_time, u.update_time, u.delete_time FROM users u WHERE u.organization_id = $1;",
		},
		{
			name:       "should list users with limit",
			ctx:        context.Background(),
			db:         database.NewDB(dbtest.Connect(context.Background(), t, dbc.ConnectionString)),
			pagination: storage.Pagination{Limit: 2, Offset: 0},
			expectedUsers: []*storage.User{
				seed.NewUser(
					"7f62b8ca-6c97-4081-adbc-2b4611b41617",
					"7bb443e5-8974-44c2-8b7c-b95124205264",
					"test",
					"active",
					"test@test.com",
					"",
					time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Time{},
				),
				seed.NewUser(
					"92c02cd8-9286-4687-bb4d-60ee95c769ed",
					"7bb443e5-8974-44c2-8b7c-b95124205264",
					"test2",
					"active",
					"test2@test.com",
					"",
					time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Time{},
				),
			},
			expectedQuery: "SELECT u.id, u.organization_id, u.display_name, u.email, u.status, u.create_time, u.update_time, u.delete_time FROM users u LIMIT $1;",
		},
		{
			name:       "should list users with offset",
			ctx:        context.Background(),
			db:         database.NewDB(dbtest.Connect(context.Background(), t, dbc.ConnectionString)),
			pagination: storage.Pagination{Limit: 0, Offset: 1},
			expectedUsers: []*storage.User{
				seed.NewUser(
					"92c02cd8-9286-4687-bb4d-60ee95c769ed",
					"7bb443e5-8974-44c2-8b7c-b95124205264",
					"test2",
					"active",
					"test2@test.com",
					"",
					time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Time{},
				),
				seed.NewUser(
					"2198aab6-2ece-429a-8d7c-1654ab8b7d8f",
					"7bb443e5-8974-44c2-8b7c-b95124205264",
					"test3",
					"active",
					"test3@test.com",
					"",
					time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Time{},
				),
			},
			expectedQuery: "SELECT u.id, u.organization_id, u.display_name, u.email, u.status, u.create_time, u.update_time, u.delete_time FROM users u OFFSET $1;",
		},
		{
			name:       "should list users with limit and offset",
			ctx:        context.Background(),
			db:         database.NewDB(dbtest.Connect(context.Background(), t, dbc.ConnectionString)),
			pagination: storage.Pagination{Limit: 1, Offset: 1},
			expectedUsers: []*storage.User{
				seed.NewUser(
					"92c02cd8-9286-4687-bb4d-60ee95c769ed",
					"7bb443e5-8974-44c2-8b7c-b95124205264",
					"test2",
					"active",
					"test2@test.com",
					"",
					time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Time{},
				),
			},
			expectedQuery: "SELECT u.id, u.organization_id, u.display_name, u.email, u.status, u.create_time, u.update_time, u.delete_time FROM users u LIMIT $1 OFFSET $2;",
		},
		{
			name: "should list users with sort",
			ctx:  context.Background(),
			db:   database.NewDB(dbtest.Connect(context.Background(), t, dbc.ConnectionString)),
			orderBy: storage.UserOrderBy{
				{Field: storage.UserDisplayName, Direction: storage.DESC},
				{Field: storage.UserCreateTime, Direction: storage.DESC},
			},
			expectedUsers: []*storage.User{
				seed.NewUser(
					"2198aab6-2ece-429a-8d7c-1654ab8b7d8f",
					"7bb443e5-8974-44c2-8b7c-b95124205264",
					"test3",
					"active",
					"test3@test.com",
					"",
					time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Time{},
				),
				seed.NewUser(
					"92c02cd8-9286-4687-bb4d-60ee95c769ed",
					"7bb443e5-8974-44c2-8b7c-b95124205264",
					"test2",
					"active",
					"test2@test.com",
					"",
					time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Time{},
				),
				seed.NewUser(
					"7f62b8ca-6c97-4081-adbc-2b4611b41617",
					"7bb443e5-8974-44c2-8b7c-b95124205264",
					"test",
					"active",
					"test@test.com",
					"",
					time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Time{},
				),
			},
			expectedQuery: "SELECT u.id, u.organization_id, u.display_name, u.email, u.status, u.create_time, u.update_time, u.delete_time FROM users u ORDER BY u.display_name DESC, u.create_time DESC;",
		},
		{
			name: "should return error if sorting field is unknown",
			ctx:  context.Background(),
			db:   database.NewDB(dbtest.Connect(context.Background(), t, dbc.ConnectionString)),
			orderBy: storage.UserOrderBy{
				{Field: storage.UserField("unknown field"), Direction: storage.DESC},
			},
			expectedError: storage.ErrUserUnknownField,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repo, recorder := connectAndCreateRepo(t, tc.db)
			var users []*storage.User
			var err error

			users, err = repo.List(tc.ctx, tc.pagination, tc.orderBy, tc.conditions...)

			assertError(t, tc.expectedError, err)
			assertUsers(t, tc.expectedUsers, users)
			assertQuery(t, recorder, tc.expectedQuery)
		})
	}
}
