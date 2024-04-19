package postgres_test

import (
	"context"
	_ "embed"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/dwethmar/lingo/cmd/auth/migrations"
	seedPostgres "github.com/dwethmar/lingo/cmd/auth/storage/seed/postgres"
	"github.com/dwethmar/lingo/cmd/auth/storage/user"
	"github.com/dwethmar/lingo/cmd/auth/storage/user/postgres"
	"github.com/dwethmar/lingo/pkg/database"
	"github.com/dwethmar/lingo/pkg/database/dbtesting"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

//go:embed testdata/update_users.sql
var updateUserSQLQuery []byte

func NewUser(
	id string,
	username string,
	email string,
	password string,
	createTime time.Time,
	updateTime time.Time,
) *user.User {
	return &user.User{
		ID:         uuid.Must(uuid.Parse(id)),
		Username:   username,
		Email:      email,
		Password:   password,
		CreateTime: createTime,
		UpdateTime: updateTime,
	}
}

// setupTestDB runs the migrations.
func dbSetup(dbURL string) error {
	return dbtesting.Migrate(dbURL, migrations.FS)
}

func TestNew(t *testing.T) {
	t.Run("should return a new repository", func(t *testing.T) {
		if got := postgres.New(nil); got == nil {
			t.Error("expected repository")
		}
	})
}

func TestRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	dbc := dbtesting.SetupTestDB(t, "auth", dbSetup)

	t.Run("Create should create a new user", func(t *testing.T) {
		ctx := context.Background()
		db := dbtesting.ConnectTestDB(ctx, t, dbc.ConnectionString)

		repo := postgres.New(db)
		user, err := repo.Create(ctx, NewUser(
			"35297169-89d8-444d-8499-c6341e3a0770",
			"test",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		))

		if err != nil {
			t.Fatal(err)
		}

		expect := NewUser(
			"35297169-89d8-444d-8499-c6341e3a0770",
			"test",
			"test@test.com",
			"",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		)

		if diff := cmp.Diff(expect, user); diff != "" {
			t.Errorf("Create() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("should return an error if the user id already exists", func(t *testing.T) {
		ctx := context.Background()
		db := dbtesting.ConnectTestDB(ctx, t, dbc.ConnectionString)

		repo := postgres.New(db)
		_, err := repo.Create(ctx, NewUser(
			"485819f0-9e48-4d25-b07b-6de8a2076be2",
			"username_300",
			"test_300@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		))

		if err != nil {
			t.Fatal(err)
		}

		_, err = repo.Create(ctx, NewUser(
			"485819f0-9e48-4d25-b07b-6de8a2076be2",
			"username_301",
			"test_3001@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		))

		if !errors.Is(err, user.ErrUniqueIDConflict) {
			t.Errorf("expected %q, got %q", user.ErrUniqueIDConflict, err)
		}
	})

	t.Run("should return an error if the user username already exists", func(t *testing.T) {
		ctx := context.Background()
		db := dbtesting.ConnectTestDB(ctx, t, dbc.ConnectionString)

		repo := postgres.New(db)
		_, err := repo.Create(ctx, NewUser(
			"4c3362d9-3956-4b15-b839-b5791460a518",
			"username",
			"test_500@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		))

		if err != nil {
			t.Fatal(err)
		}

		_, err = repo.Create(ctx, NewUser(
			"8f56d098-731f-48a3-ab19-942f7c793732",
			"username",
			"test_501@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		))

		if !errors.Is(err, user.ErrUniqueUsernameConflict) {
			t.Errorf("expected %q, got %q", user.ErrUniqueUsernameConflict, err)
		}
	})

	t.Run("should return an error if the user email already exists", func(t *testing.T) {
		ctx := context.Background()
		db := dbtesting.ConnectTestDB(ctx, t, dbc.ConnectionString)

		repo := postgres.New(db)
		_, err := repo.Create(ctx, NewUser(
			"2e56b481-05fe-4ce3-b072-a94fbf8aeab3",
			"username_400",
			"test_400@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		))

		if err != nil {
			t.Fatal(err)
		}

		_, err = repo.Create(ctx, NewUser(
			"5e6f2f35-1de1-4803-8fdd-9b67706f887e",
			"username_401",
			"test_400@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		))

		if !errors.Is(err, user.ErrUniqueEmailConflict) {
			t.Errorf("expected %q, got %q", user.ErrUniqueEmailConflict, err)
		}
	})
}

func TestRepository_Get(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	dbc := dbtesting.SetupTestDB(t, "auth", dbSetup)

	t.Run("Get should get a user", func(t *testing.T) {
		ctx := context.Background()
		db := dbtesting.ConnectTestDB(ctx, t, dbc.ConnectionString)

		if err := seedPostgres.User(ctx, db, NewUser(
			"35297169-89d8-444d-8499-c6341e3a0770",
			"test",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		)); err != nil {
			t.Fatal(err)
		}

		repo := postgres.New(db)
		user, err := repo.Get(ctx, uuid.Must(uuid.Parse("35297169-89d8-444d-8499-c6341e3a0770")))

		if err != nil {
			t.Fatal(err)
		}

		expect := NewUser(
			"35297169-89d8-444d-8499-c6341e3a0770",
			"test",
			"test@test.com",
			"",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		)

		if diff := cmp.Diff(expect, user); diff != "" {
			t.Errorf("Create() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("should return an error if the user does not exist", func(t *testing.T) {
		ctx := context.Background()
		db := dbtesting.ConnectTestDB(ctx, t, dbc.ConnectionString)

		repo := postgres.New(db)
		u, err := repo.Get(ctx, uuid.Must(uuid.Parse("946adb15-195e-44df-922b-4a45b9505684")))

		if !errors.Is(err, user.ErrNotFound) {
			t.Errorf("expected %q, got %q", user.ErrNotFound, err)
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

	dbc := dbtesting.SetupTestDB(t, "auth", dbSetup)

	t.Run("should update a user", func(t *testing.T) {
		ctx := context.Background()
		db := dbtesting.ConnectTestDB(ctx, t, dbc.ConnectionString)

		if err := seedPostgres.User(ctx, db, NewUser(
			"957b12c5-1071-40d9-8bec-6ed195c8cfbf",
			"test",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		)); err != nil {
			t.Fatal(err)
		}

		recorder := database.NewRecorder(db)
		repo := postgres.New(recorder)

		user, err := repo.Update(ctx, NewUser(
			"957b12c5-1071-40d9-8bec-6ed195c8cfbf",
			"test2", // updated username
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		), user.Username, user.Email, user.Password, user.UpdateTime)

		if err != nil {
			t.Fatal(err)
		}

		expect := NewUser(
			"957b12c5-1071-40d9-8bec-6ed195c8cfbf",
			"test2",
			"test@test.com",
			"",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		)

		if diff := cmp.Diff(expect, user); diff != "" {
			t.Errorf("Create() mismatch (-want +got):\n%s", diff)
		}

		// check if only the username was updated
		expectedQuery := strings.TrimSpace(strings.ReplaceAll(string(updateUserSQLQuery), "\n", " "))
		gotQuery := recorder.RowQueries[0].Query
		gotQuery = strings.TrimSpace(strings.ReplaceAll(gotQuery, "\n", " "))

		if gotQuery != expectedQuery {
			t.Errorf("expected %q, got %q", expectedQuery, gotQuery)
		}
	})

	t.Run("should return an error if the user does not exist", func(t *testing.T) {
		ctx := context.Background()
		db := dbtesting.ConnectTestDB(ctx, t, dbc.ConnectionString)

		repo := postgres.New(db)
		_, err := repo.Update(ctx, NewUser(
			"f2e8b3cd-07a3-4d7c-9eef-cf02452d8332",
			"test",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		), user.Username, user.UpdateTime)

		if !errors.Is(err, user.ErrNotFound) {
			t.Errorf("expected %q, got %q", user.ErrNotFound, err)
		}
	})

	t.Run("should return an error if no fields are provided", func(t *testing.T) {
		ctx := context.Background()
		r, err := postgres.New(nil).Update(ctx, NewUser(
			"957b12c5-1071-40d9-8bec-6ed195c8cfbf",
			"test",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		))

		if err == nil || !errors.Is(err, user.ErrNoFieldsToUpdate) {
			t.Errorf("expected %q, got %q", user.ErrNoFieldsToUpdate, err)
		}

		if r != nil {
			t.Errorf("expected nil, got %v", r)
		}
	})
}

func TestRepository_GetByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	dbc := dbtesting.SetupTestDB(t, "auth", dbSetup)

	t.Run("GetByUsername should get a user by username", func(t *testing.T) {
		ctx := context.Background()
		db := dbtesting.ConnectTestDB(ctx, t, dbc.ConnectionString)

		if err := seedPostgres.User(ctx, db, NewUser(
			"82651da9-c2ff-4152-8eae-7555d5a42aad",
			"test",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		)); err != nil {
			t.Fatal(err)
		}

		repo := postgres.New(db)
		user, err := repo.GetByEmail(ctx, "test@test.com")

		if err != nil {
			t.Fatal(err)
		}

		expect := NewUser(
			"82651da9-c2ff-4152-8eae-7555d5a42aad",
			"test",
			"test@test.com",
			"",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		)

		if diff := cmp.Diff(expect, user); diff != "" {
			t.Errorf("Create() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("should return an error if the user does not exist", func(t *testing.T) {
		ctx := context.Background()
		db := dbtesting.ConnectTestDB(ctx, t, dbc.ConnectionString)

		repo := postgres.New(db)
		u, err := repo.GetByEmail(ctx, "test2@test.com")

		if err == nil || !errors.Is(err, user.ErrNotFound) {
			t.Errorf("expected %q, got %q", user.ErrNotFound, err)
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

	dbc := dbtesting.SetupTestDB(t, "auth", dbSetup)

	t.Run("Delete should delete a user", func(t *testing.T) {
		ctx := context.Background()
		db := dbtesting.ConnectTestDB(ctx, t, dbc.ConnectionString)

		if err := seedPostgres.User(ctx, db, NewUser(
			"82651da9-c2ff-4152-8eae-7555d5a42aad",
			"test",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		)); err != nil {
			t.Fatal(err)
		}

		repo := postgres.New(db)
		if err := repo.Delete(ctx, uuid.Must(uuid.Parse("82651da9-c2ff-4152-8eae-7555d5a42aad"))); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should return an error if the user does not exist", func(t *testing.T) {
		ctx := context.Background()
		db := dbtesting.ConnectTestDB(ctx, t, dbc.ConnectionString)

		repo := postgres.New(db)
		err := repo.Delete(ctx, uuid.Must(uuid.Parse("82651da9-c2ff-4152-8eae-7555d5a42aad")))

		if !errors.Is(err, user.ErrNotFound) {
			t.Errorf("expected %q, got %q", user.ErrNotFound, err)
		}
	})
}
