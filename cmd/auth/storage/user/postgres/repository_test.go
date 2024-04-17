package postgres

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/dwethmar/lingo/cmd/auth/migrations"
	"github.com/dwethmar/lingo/cmd/auth/storage/user"
	"github.com/dwethmar/lingo/pkg/database"
	"github.com/dwethmar/lingo/pkg/database/dbtesting"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

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

func TestNewRepository(t *testing.T) {
	t.Run("should return a new repository", func(t *testing.T) {
		if got := NewRepository(nil); got == nil {
			t.Error("expected repository")
		}
	})
}

func TestRepository_Create(t *testing.T) {
	ctx := context.Background()
	dbc, err := dbtesting.SetupPostgresContainer(context.Background(), func(dbURL string) error {
		return dbtesting.Migrate(dbURL, migrations.FS)
	})

	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := dbc.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	t.Run("Create should create a new user", func(t *testing.T) {
		t.Cleanup(func() {
			err = dbc.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		db, close := dbtesting.SetupTestDB(ctx, t, dbc)
		defer close()

		repo := NewRepository(db)
		user, err := repo.Create(context.Background(), NewUser(
			"485819f0-9e48-4d25-b07b-6de8a2076be2",
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
			"485819f0-9e48-4d25-b07b-6de8a2076be2",
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
		t.Cleanup(func() {
			err = dbc.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		db, close := dbtesting.SetupTestDB(ctx, t, dbc)
		defer close()

		repo := NewRepository(db)
		_, err = repo.Create(context.Background(), NewUser(
			"485819f0-9e48-4d25-b07b-6de8a2076be2",
			"test",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		))

		if err != nil {
			t.Fatal(err)
		}

		_, err = repo.Create(context.Background(), NewUser(
			"485819f0-9e48-4d25-b07b-6de8a2076be2",
			"test",
			"test@test.com",
			"password",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		))

		if err.Error() != "failed to insert user: pq: duplicate key value violates unique constraint \"users_pkey\"" {
			t.Error("expected error")
		}
	})
}

func TestRepository_Get(t *testing.T) {
	type fields struct {
		db database.DB
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *user.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				db: tt.fields.db,
			}
			got, err := r.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_Update(t *testing.T) {
	type fields struct {
		db database.DB
	}
	type args struct {
		ctx    context.Context
		u      *user.User
		fields []user.Field
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *user.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				db: tt.fields.db,
			}
			got, err := r.Update(tt.args.ctx, tt.args.u, tt.args.fields...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	type fields struct {
		db database.DB
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				db: tt.fields.db,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Repository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
