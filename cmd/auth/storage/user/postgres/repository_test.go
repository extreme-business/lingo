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

func TestNewRepository(t *testing.T) {
	type args struct {
		db database.DB
	}
	tests := []struct {
		name string
		args args
		want *Repository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_Create(t *testing.T) {
	ctx := context.Background()
	dbc, err := dbtesting.SetupPostgres(context.Background(), func(dbURL string) error {
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
			if err = dbc.Restore(context.Background()); err != nil {
				t.Fatal(err)
			}
		})

		expect := &user.User{
			ID:         uuid.Must(uuid.Parse("485819f0-9e48-4d25-b07b-6de8a2076be2")),
			Username:   "test",
			Email:      "wow@test.nl",
			Password:   "",
			CreateTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		db, close, err := database.Connect(ctx, dbc.URL)
		if err != nil {
			t.Fatal(err)
		}
		defer close()

		repo := NewRepository(db)
		user, err := repo.Create(context.Background(), &user.User{
			ID:         uuid.Must(uuid.Parse("485819f0-9e48-4d25-b07b-6de8a2076be2")),
			Username:   "test",
			Email:      "wow@test.nl",
			Password:   "password",
			CreateTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		})

		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(expect, user); diff != "" {
			t.Errorf("Create() mismatch (-want +got):\n%s", diff)
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
