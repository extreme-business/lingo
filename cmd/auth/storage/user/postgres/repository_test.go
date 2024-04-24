package postgres_test

import (
	"context"
	_ "embed"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/dwethmar/lingo/cmd/auth/migrations"
	"github.com/dwethmar/lingo/cmd/auth/storage/organization"
	seedPostgres "github.com/dwethmar/lingo/cmd/auth/storage/seed/postgres"
	"github.com/dwethmar/lingo/cmd/auth/storage/user"
	"github.com/dwethmar/lingo/cmd/auth/storage/user/postgres"
	"github.com/dwethmar/lingo/pkg/database/dbtest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func NewOrganization(id string, displayName string, createTime time.Time, updateTime time.Time) *organization.Organization {
	return &organization.Organization{
		ID:          uuid.Must(uuid.Parse(id)),
		DisplayName: displayName,
		CreateTime:  createTime,
		UpdateTime:  updateTime,
	}
}

func NewUser(id string, organizationID string, displayName string, email string, password string, createTime time.Time, updateTime time.Time) *user.User {
	return &user.User{
		ID:             uuid.Must(uuid.Parse(id)),
		OrganizationID: uuid.Must(uuid.Parse(organizationID)),
		DisplayName:    displayName,
		Email:          email,
		Password:       password,
		CreateTime:     createTime,
		UpdateTime:     updateTime,
	}
}

// setupTestDB runs the migrations.
func dbSetup(dbURL string) error { return dbtest.Migrate(dbURL, migrations.FS) }

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

	dbc := dbtest.SetupTestDB(t, "auth", dbSetup)
	seedPostgres.Run(t, dbc.ConnectionString, []*organization.Organization{
		NewOrganization(
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		),
	}, []*user.User{})

	t.Run("Create should create a new user", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.ConnectTestDB(ctx, t, dbc.ConnectionString)

		repo := postgres.New(db)
		user, err := repo.Create(ctx, NewUser(
			"35297169-89d8-444d-8499-c6341e3a0770",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
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
			"7bb443e5-8974-44c2-8b7c-b95124205264",
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
		db := dbtest.ConnectTestDB(ctx, t, dbc.ConnectionString)

		repo := postgres.New(db)
		_, err := repo.Create(ctx, NewUser(
			"485819f0-9e48-4d25-b07b-6de8a2076be2",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
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
			"7bb443e5-8974-44c2-8b7c-b95124205264",
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

	t.Run("should return an error if the user email already exists", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.ConnectTestDB(ctx, t, dbc.ConnectionString)

		repo := postgres.New(db)
		_, err := repo.Create(ctx, NewUser(
			"2e56b481-05fe-4ce3-b072-a94fbf8aeab3",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
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
			"7bb443e5-8974-44c2-8b7c-b95124205264",
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

	dbc := dbtest.SetupTestDB(t, "auth", dbSetup)
	seedPostgres.Run(t, dbc.ConnectionString, []*organization.Organization{
		NewOrganization(
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		),
	}, []*user.User{
		NewUser(
			"35297169-89d8-444d-8499-c6341e3a0770",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		),
	})

	t.Run("Get should get a user", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.ConnectTestDB(ctx, t, dbc.ConnectionString)

		repo := postgres.New(db)
		user, err := repo.Get(ctx, uuid.Must(uuid.Parse("35297169-89d8-444d-8499-c6341e3a0770")))

		if err != nil {
			t.Fatal(err)
		}

		expect := NewUser(
			"35297169-89d8-444d-8499-c6341e3a0770",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
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
		db := dbtest.ConnectTestDB(ctx, t, dbc.ConnectionString)
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

	dbc := dbtest.SetupTestDB(t, "auth", dbSetup)
	seedPostgres.Run(t, dbc.ConnectionString, []*organization.Organization{
		NewOrganization(
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		),
		NewOrganization(
			"f226487d-61ff-4a18-a2d9-ab888b22dbc8",
			"test2",
			time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
		),
	}, []*user.User{
		NewUser(
			"35297169-89d8-444d-8499-c6341e3a0770",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		),
	})

	t.Run("should update a user", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.ConnectTestDB(ctx, t, dbc.ConnectionString)
		recorder := dbtest.NewRecorder(db)
		repo := postgres.New(recorder)

		user, err := repo.Update(ctx, NewUser(
			"35297169-89d8-444d-8499-c6341e3a0770",
			"f226487d-61ff-4a18-a2d9-ab888b22dbc8",
			"test2", // updated username
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		), user.OrganizationID, user.DisplayName, user.Email, user.Password, user.UpdateTime)

		if err != nil {
			t.Fatal(err)
		}

		expect := NewUser(
			"35297169-89d8-444d-8499-c6341e3a0770",
			"f226487d-61ff-4a18-a2d9-ab888b22dbc8",
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
		query := recorder.RowQueries[0].Query
		query = strings.TrimSpace(strings.ReplaceAll(query, "\n", " "))
		expectedQuery := "UPDATE users SET organization_id = $1, display_name = $2, email = $3, password = $4, update_time = $5 WHERE id = $6 RETURNING id, organization_id, display_name, email, create_time, update_time;"

		if query != expectedQuery {
			t.Errorf("expected %q, got %q", expectedQuery, query)
		}
	})

	t.Run("should return an error if the user does not exist", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.ConnectTestDB(ctx, t, dbc.ConnectionString)

		repo := postgres.New(db)
		_, err := repo.Update(ctx, NewUser(
			"f2e8b3cd-07a3-4d7c-9eef-cf02452d8332",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		), user.DisplayName, user.UpdateTime)

		if !errors.Is(err, user.ErrNotFound) {
			t.Errorf("expected %q, got %q", user.ErrNotFound, err)
		}
	})

	t.Run("should return an error if no fields are provided", func(t *testing.T) {
		ctx := context.Background()
		r, err := postgres.New(nil).Update(ctx, NewUser(
			"957b12c5-1071-40d9-8bec-6ed195c8cfbf",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
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

func TestRepository_Update_fields(t *testing.T) {
	dbc := dbtest.SetupTestDB(t, "auth", dbSetup)
	seedPostgres.Run(t, dbc.ConnectionString, []*organization.Organization{
		NewOrganization(
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		),
	}, []*user.User{
		NewUser(
			"957b12c5-1071-40d9-8bec-6ed195c8cfbf",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		),
	})

	tests := []struct {
		name   string
		user   *user.User
		fields []user.Field
		err    error
	}{
		{
			name: "should return an error no fields are provided",
			user: NewUser(
				"957b12c5-1071-40d9-8bec-6ed195c8cfbf",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test@test.com",
				"password",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			),
			fields: []user.Field{},
			err:    user.ErrNoFieldsToUpdate,
		},
		{
			name: "should return an error field is create_time",
			user: NewUser(
				"957b12c5-1071-40d9-8bec-6ed195c8cfbf",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test@test.com",
				"password",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			),
			fields: []user.Field{user.CreateTime},
			err:    postgres.ErrUpdateReadOnlyCreateTime,
		},
		{
			name: "should return an error if the field is unknown",
			user: NewUser(
				"957b12c5-1071-40d9-8bec-6ed195c8cfbf",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test@test.com",
				"password",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			),
			fields: []user.Field{"unknown"},
			err:    postgres.ErrUpdateUnknownField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := postgres.New(nil).Update(ctx, tt.user, tt.fields...)

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

	dbc := dbtest.SetupTestDB(t, "auth", dbSetup)
	seedPostgres.Run(t, dbc.ConnectionString, []*organization.Organization{
		NewOrganization(
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		),
	}, []*user.User{
		NewUser(
			"82651da9-c2ff-4152-8eae-7555d5a42aad",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		),
	})

	t.Run("should get a user by email", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.ConnectTestDB(ctx, t, dbc.ConnectionString)
		repo := postgres.New(db)
		user, err := repo.GetByEmail(ctx, "test@test.com")

		if err != nil {
			t.Fatal(err)
		}

		expect := NewUser(
			"82651da9-c2ff-4152-8eae-7555d5a42aad",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		)

		if diff := cmp.Diff(expect, user); diff != "" {
			t.Errorf("Create() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("should return an error if the user does not exist", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.ConnectTestDB(ctx, t, dbc.ConnectionString)

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

	dbc := dbtest.SetupTestDB(t, "auth", dbSetup)
	seedPostgres.Run(t, dbc.ConnectionString, []*organization.Organization{
		NewOrganization(
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		),
	}, []*user.User{
		NewUser(
			"35297169-89d8-444d-8499-c6341e3a0770",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		),
	})

	t.Run("Delete should delete a user", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.ConnectTestDB(ctx, t, dbc.ConnectionString)
		repo := postgres.New(db)
		if err := repo.Delete(ctx, uuid.Must(uuid.Parse("35297169-89d8-444d-8499-c6341e3a0770"))); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should return an error if the user does not exist", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.ConnectTestDB(ctx, t, dbc.ConnectionString)

		repo := postgres.New(db)
		err := repo.Delete(ctx, uuid.Must(uuid.Parse("82651da9-c2ff-4152-8eae-7555d5a42aad")))

		if !errors.Is(err, user.ErrNotFound) {
			t.Errorf("expected %q, got %q", user.ErrNotFound, err)
		}
	})
}

func TestRepository_List(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	dbc := dbtest.SetupTestDB(t, "auth", dbSetup)
	seedPostgres.Run(t, dbc.ConnectionString, []*organization.Organization{
		NewOrganization(
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		),
	}, []*user.User{
		NewUser(
			"7f62b8ca-6c97-4081-adbc-2b4611b41617",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			"test1@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
		),
		NewUser(
			"92c02cd8-9286-4687-bb4d-60ee95c769ed",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			"test2@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC),
		),
		NewUser(
			"2198aab6-2ece-429a-8d7c-1654ab8b7d8f",
			"7bb443e5-8974-44c2-8b7c-b95124205264",
			"test",
			"test3@example.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC),
		),
	})

	t.Run("should list users", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.ConnectTestDB(ctx, t, dbc.ConnectionString)
		recorder := dbtest.NewRecorder(db)
		repo := postgres.New(recorder)
		users, err := repo.List(ctx, user.Pagination{Limit: 0, Offset: 0}, []user.Sort{})

		if err != nil {
			t.Fatal(err)
		}

		expect := []*user.User{
			NewUser(
				"7f62b8ca-6c97-4081-adbc-2b4611b41617",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test1@test.com",
				"",
				time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
			),
			NewUser(
				"92c02cd8-9286-4687-bb4d-60ee95c769ed",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test2@test.com",
				"",
				time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC),
			),
			NewUser(
				"2198aab6-2ece-429a-8d7c-1654ab8b7d8f",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test3@example.com",
				"",
				time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC),
			),
		}

		if diff := cmp.Diff(expect, users); diff != "" {
			t.Errorf("Create() mismatch (-want +got):\n%s", diff)
		}

		query := recorder.Queries[0].Query
		query = strings.TrimSpace(strings.ReplaceAll(query, "\n", " "))
		expectedQuery := "SELECT id, organization_id, display_name, email, create_time, update_time FROM users;"

		if query != expectedQuery {
			t.Errorf("expected %q, got %q", expectedQuery, query)
		}
	})

	t.Run("should list users with limit", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.ConnectTestDB(ctx, t, dbc.ConnectionString)
		recorder := dbtest.NewRecorder(db)
		repo := postgres.New(recorder)
		users, err := repo.List(ctx, user.Pagination{Limit: 2, Offset: 0}, []user.Sort{})

		if err != nil {
			t.Fatal(err)
		}

		expect := []*user.User{
			NewUser(
				"7f62b8ca-6c97-4081-adbc-2b4611b41617",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test1@test.com",
				"",
				time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
			),
			NewUser(
				"92c02cd8-9286-4687-bb4d-60ee95c769ed",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test2@test.com",
				"",
				time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC),
			),
		}

		if diff := cmp.Diff(expect, users); diff != "" {
			t.Errorf("Create() mismatch (-want +got):\n%s", diff)
		}

		query := recorder.Queries[0].Query
		query = strings.TrimSpace(strings.ReplaceAll(query, "\n", " "))
		expectedQuery := "SELECT id, organization_id, display_name, email, create_time, update_time FROM users LIMIT $1;"

		if query != expectedQuery {
			t.Errorf("expected %q, got %q", expectedQuery, query)
		}
	})

	t.Run("should list users with offset", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.ConnectTestDB(ctx, t, dbc.ConnectionString)
		recorder := dbtest.NewRecorder(db)
		repo := postgres.New(recorder)
		users, err := repo.List(ctx, user.Pagination{Limit: 0, Offset: 1}, []user.Sort{})

		if err != nil {
			t.Fatal(err)
		}

		expect := []*user.User{
			NewUser(
				"92c02cd8-9286-4687-bb4d-60ee95c769ed",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test2@test.com",
				"",
				time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC),
			),
			NewUser(
				"2198aab6-2ece-429a-8d7c-1654ab8b7d8f",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test3@example.com",
				"",
				time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC),
			),
		}

		if diff := cmp.Diff(expect, users); diff != "" {
			t.Errorf("Create() mismatch (-want +got):\n%s", diff)
		}

		query := recorder.Queries[0].Query
		query = strings.TrimSpace(strings.ReplaceAll(query, "\n", " "))
		expectedQuery := "SELECT id, organization_id, display_name, email, create_time, update_time FROM users OFFSET $1;"

		if query != expectedQuery {
			t.Errorf("expected %q, got %q", expectedQuery, query)
		}
	})

	t.Run("should list users with limit and offset", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.ConnectTestDB(ctx, t, dbc.ConnectionString)
		recorder := dbtest.NewRecorder(db)
		repo := postgres.New(recorder)
		users, err := repo.List(ctx, user.Pagination{Limit: 1, Offset: 1}, []user.Sort{})

		if err != nil {
			t.Fatal(err)
		}

		expect := []*user.User{
			NewUser(
				"92c02cd8-9286-4687-bb4d-60ee95c769ed",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test2@test.com",
				"",
				time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC),
			),
		}

		if diff := cmp.Diff(expect, users); diff != "" {
			t.Errorf("Create() mismatch (-want +got):\n%s", diff)
		}

		query := recorder.Queries[0].Query
		query = strings.TrimSpace(strings.ReplaceAll(query, "\n", " "))
		expectedQuery := "SELECT id, organization_id, display_name, email, create_time, update_time FROM users LIMIT $1 OFFSET $2;"

		if query != expectedQuery {
			t.Errorf("expected %q, got %q", expectedQuery, query)
		}
	})

	t.Run("should list users with sort", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.ConnectTestDB(ctx, t, dbc.ConnectionString)
		recorder := dbtest.NewRecorder(db)
		repo := postgres.New(recorder)
		users, err := repo.List(ctx, user.Pagination{}, []user.Sort{
			{Field: user.DisplayName, Direction: user.DESC},
			{Field: user.CreateTime, Direction: user.DESC},
		})

		if err != nil {
			t.Fatal(err)
		}

		expect := []*user.User{
			NewUser(
				"2198aab6-2ece-429a-8d7c-1654ab8b7d8f",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test3@example.com",
				"",
				time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC),
			),
			NewUser(
				"92c02cd8-9286-4687-bb4d-60ee95c769ed",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test2@test.com",
				"",
				time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC),
			),
			NewUser(
				"7f62b8ca-6c97-4081-adbc-2b4611b41617",
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				"test1@test.com",
				"",
				time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
			),
		}

		if diff := cmp.Diff(expect, users); diff != "" {
			t.Errorf("Create() mismatch (-want +got):\n%s", diff)
		}

		query := recorder.Queries[0].Query
		query = strings.TrimSpace(strings.ReplaceAll(query, "\n", " "))
		expectedQuery := "SELECT id, organization_id, display_name, email, create_time, update_time FROM users ORDER BY display_name DESC, create_time DESC;"

		if query != expectedQuery {
			t.Errorf("expected %q, got %q", expectedQuery, query)
		}
	})

	t.Run("should return error if sorting field is unknown", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.ConnectTestDB(ctx, t, dbc.ConnectionString)
		recorder := dbtest.NewRecorder(db)
		repo := postgres.New(recorder)
		users, err := repo.List(ctx, user.Pagination{}, []user.Sort{
			{Field: user.Field("unknown field"), Direction: user.DESC},
		})

		if err == nil || !errors.Is(err, user.ErrInvalidSortField) {
			t.Errorf("expected %q, got %q", user.ErrInvalidSortField, err)
		}

		if users != nil {
			t.Errorf("expected nil, got %v", users)
		}

		if len(recorder.Queries) > 0 {
			t.Error("expected no queries")
		}
	})
}
