package domain

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/extreme-business/lingo/apps/account/storage"
	protoaccount "github.com/extreme-business/lingo/proto/gen/go/public/account/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	userStorageMock "github.com/extreme-business/lingo/apps/account/storage/mock/user"
)

func TestUserStatus_String(t *testing.T) {
	tests := []struct {
		name string
		s    UserStatus
		want string
	}{
		{
			name: "UserStatusActive",
			s:    UserStatusActive,
			want: "active",
		},
		{
			name: "UserStatusInactive",
			s:    UserStatusInactive,
			want: "inactive",
		},
		{
			name: "UserStatusDeleted",
			s:    UserStatusDeleted,
			want: "deleted",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("UserStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_FromProto(t *testing.T) {
	t.Run("when the user is valid", func(t *testing.T) {
		in := &protoaccount.User{
			Name:        "organizations/99f85dd9-6df3-4320-8acd-d665a459be38/users/18cdf53a-4c02-46df-869a-ba12f2e90d35",
			DisplayName: "test",
			Email:       "test@test.com",
			Password:    "hashedpassword",
			CreateTime:  timestamppb.New(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			UpdateTime:  timestamppb.New(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
			DeleteTime:  timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		}

		u := &User{}
		if err := u.FromProto(in); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := &User{
			ID:             uuid.MustParse("18cdf53a-4c02-46df-869a-ba12f2e90d35"),
			OrganizationID: uuid.MustParse("99f85dd9-6df3-4320-8acd-d665a459be38"),
			DisplayName:    "test",
			Email:          "test@test.com",
			HashedPassword: "",
			CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			DeleteTime:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		if diff := cmp.Diff(u, expected); diff != "" {
			t.Fatalf("unexpected diff: %v", diff)
		}
	})
}

func TestUser_ToProto(t *testing.T) {
	t.Run("when the user is valid", func(t *testing.T) {
		u := &User{
			ID:             uuid.MustParse("18cdf53a-4c02-46df-869a-ba12f2e90d35"),
			OrganizationID: uuid.MustParse("99f85dd9-6df3-4320-8acd-d665a459be38"),
			DisplayName:    "test",
			Email:          "test@test.com",
			HashedPassword: "",
			CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			DeleteTime:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		out := &protoaccount.User{}
		if err := u.ToProto(out); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := &protoaccount.User{
			Name:        "organizations/99f85dd9-6df3-4320-8acd-d665a459be38/users/18cdf53a-4c02-46df-869a-ba12f2e90d35",
			DisplayName: "test",
			Email:       "test@test.com",
			Password:    "",
			CreateTime:  timestamppb.New(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			UpdateTime:  timestamppb.New(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
			DeleteTime:  timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		}

		// Use cmpopts.IgnoreUnexported to handle protobuf's internal fields
		if diff := cmp.Diff(expected, out, cmpopts.IgnoreUnexported(protoaccount.User{}, timestamppb.Timestamp{})); diff != "" {
			t.Fatalf("unexpected diff: %v", diff)
		}
	})
}

func TestUser_ToStorage(t *testing.T) {
	t.Run("when the user is valid", func(t *testing.T) {
		u := &User{
			ID:             uuid.MustParse("18cdf53a-4c02-46df-869a-ba12f2e90d35"),
			OrganizationID: uuid.MustParse("99f85dd9-6df3-4320-8acd-d665a459be38"),
			DisplayName:    "test",
			Email:          "test@test.com",
			HashedPassword: "hashedpassword",
			Status:         UserStatusActive,
			CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			DeleteTime:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		out := &storage.User{}
		if err := u.ToStorage(out); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := &storage.User{
			ID:             uuid.MustParse("18cdf53a-4c02-46df-869a-ba12f2e90d35"),
			OrganizationID: uuid.MustParse("99f85dd9-6df3-4320-8acd-d665a459be38"),
			DisplayName:    "test",
			Email:          "test@test.com",
			HashedPassword: "hashedpassword",
			Status:         "active",
			CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			DeleteTime:     sql.NullTime{Time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true},
		}

		if diff := cmp.Diff(out, expected); diff != "" {
			t.Fatalf("unexpected diff: %v", diff)
		}
	})
}

func TestUser_FromStorage(t *testing.T) {
	t.Run("when the user is valid", func(t *testing.T) {
		in := &storage.User{
			ID:             uuid.MustParse("18cdf53a-4c02-46df-869a-ba12f2e90d35"),
			OrganizationID: uuid.MustParse("99f85dd9-6df3-4320-8acd-d665a459be38"),
			DisplayName:    "test",
			Email:          "test@test.com",
			HashedPassword: "hashedpassword",
			Status:         "active",
			CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			DeleteTime:     sql.NullTime{Time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true},
		}

		u := &User{}
		if err := u.FromStorage(in); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := &User{
			ID:             uuid.MustParse("18cdf53a-4c02-46df-869a-ba12f2e90d35"),
			OrganizationID: uuid.MustParse("99f85dd9-6df3-4320-8acd-d665a459be38"),
			DisplayName:    "test",
			Email:          "test@test.com",
			HashedPassword: "hashedpassword",
			Status:         UserStatusActive,
			CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			DeleteTime:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		if diff := cmp.Diff(u, expected); diff != "" {
			t.Fatalf("unexpected diff: %v", diff)
		}
	})
}

func TestNewUserReader(t *testing.T) {
	t.Run("should return non nil", func(t *testing.T) {
		r := NewUserReader(nil)
		if r == nil {
			t.Fatalf("unexpected nil")
		}
	})
}

func TestUserReader_Get(t *testing.T) {
	t.Run("when the user is found", func(t *testing.T) {
		ctx := context.Background()

		in := &storage.User{
			ID:             uuid.MustParse("18cdf53a-4c02-46df-869a-ba12f2e90d35"),
			OrganizationID: uuid.MustParse("99f85dd9-6df3-4320-8acd-d665a459be38"),
			DisplayName:    "test",
			Email:          "test@test.com",
			HashedPassword: "hashedpassword",
			Status:         "active",
			CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			DeleteTime:     sql.NullTime{Time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true},
		}

		m := &userStorageMock.Repository{
			GetFunc: func(ctx context.Context, id uuid.UUID) (*storage.User, error) {
				if id != uuid.MustParse("18cdf53a-4c02-46df-869a-ba12f2e90d35") {
					return nil, errors.New("not found")
				}
				return in, nil
			},
		}

		r := NewUserReader(m)
		out, err := r.Get(ctx, uuid.MustParse("18cdf53a-4c02-46df-869a-ba12f2e90d35"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := &User{
			ID:             uuid.MustParse("18cdf53a-4c02-46df-869a-ba12f2e90d35"),
			OrganizationID: uuid.MustParse("99f85dd9-6df3-4320-8acd-d665a459be38"),
			DisplayName:    "test",
			Email:          "test@test.com",
			HashedPassword: "hashedpassword",
			Status:         UserStatusActive,
			CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			DeleteTime:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		if diff := cmp.Diff(out, expected); diff != "" {
			t.Fatalf("unexpected diff: %v", diff)
		}
	})
}

func TestUserReader_List(t *testing.T) {
	t.Run("when the users are found", func(t *testing.T) {
		ctx := context.Background()
		in := []*storage.User{
			{
				ID:             uuid.MustParse("18cdf53a-4c02-46df-869a-ba12f2e90d35"),
				OrganizationID: uuid.MustParse("99f85dd9-6df3-4320-8acd-d665a459be38"),
				DisplayName:    "test1",
				Email:          "test1@test.com",
				HashedPassword: "hashedpassword",
				Status:         "active",
				CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdateTime:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:             uuid.MustParse("3135c0a5-f34d-4f08-8686-98193b33e743"),
				OrganizationID: uuid.MustParse("dfd8bf3e-880e-4f23-8b4e-3520ee264bf9"),
				DisplayName:    "test2",
				Email:          "test2@test.com",
				HashedPassword: "hashedpassword",
				Status:         "active",
				CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdateTime:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:             uuid.MustParse("59122a07-6068-45c4-b112-8d379ec6ea68"),
				OrganizationID: uuid.MustParse("61009b85-4eee-4f18-8fde-5ee7dbae1cc8"),
				DisplayName:    "test3",
				Email:          "test3@test.com",
				HashedPassword: "hashedpassword",
				Status:         "active",
				CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdateTime:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}

		m := &userStorageMock.Repository{
			ListFunc: func(ctx context.Context, p storage.Pagination, s storage.UserOrderBy, c ...storage.Condition) ([]*storage.User, error) {
				return in, nil
			},
		}

		r := NewUserReader(m)
		out, err := r.List(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := []*User{
			{
				ID:             uuid.MustParse("18cdf53a-4c02-46df-869a-ba12f2e90d35"),
				OrganizationID: uuid.MustParse("99f85dd9-6df3-4320-8acd-d665a459be38"),
				DisplayName:    "test1",
				Email:          "test1@test.com",
				HashedPassword: "hashedpassword",
				Status:         UserStatusActive,
				CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdateTime:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:             uuid.MustParse("3135c0a5-f34d-4f08-8686-98193b33e743"),
				OrganizationID: uuid.MustParse("dfd8bf3e-880e-4f23-8b4e-3520ee264bf9"),
				DisplayName:    "test2",
				Email:          "test2@test.com",
				HashedPassword: "hashedpassword",
				Status:         UserStatusActive,
				CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdateTime:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:             uuid.MustParse("59122a07-6068-45c4-b112-8d379ec6ea68"),
				OrganizationID: uuid.MustParse("61009b85-4eee-4f18-8fde-5ee7dbae1cc8"),
				DisplayName:    "test3",
				Email:          "test3@test.com",
				HashedPassword: "hashedpassword",
				Status:         UserStatusActive,
				CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdateTime:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}

		if diff := cmp.Diff(out, expected); diff != "" {
			t.Fatalf("unexpected diff: %v", diff)
		}
	})
}

func TestNewUserWriter(t *testing.T) {
	t.Run("should return non nil", func(t *testing.T) {
		w := NewUserWriter(nil, nil)
		if w == nil {
			t.Fatalf("unexpected nil")
		}
	})
}

func TestUserWriter_Create(t *testing.T) {
	t.Run("when the user is valid", func(t *testing.T) {
		ctx := context.Background()

		in := &User{
			ID:             uuid.MustParse("18cdf53a-4c02-46df-869a-ba12f2e90d35"),
			OrganizationID: uuid.MustParse("99f85dd9-6df3-4320-8acd-d665a459be38"),
			DisplayName:    "test",
			Email:          "test@test.com",
			HashedPassword: "hashedpassword",
			Status:         UserStatusActive,
			CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		m := &userStorageMock.Repository{
			CreateFunc: func(ctx context.Context, u *storage.User) (*storage.User, error) {
				return u, nil
			},
		}

		w := NewUserWriter(func() time.Time { return time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC) }, m)
		out, err := w.Create(ctx, in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := &User{
			ID:             uuid.MustParse("18cdf53a-4c02-46df-869a-ba12f2e90d35"),
			OrganizationID: uuid.MustParse("99f85dd9-6df3-4320-8acd-d665a459be38"),
			DisplayName:    "test",
			Email:          "test@test.com",
			HashedPassword: "hashedpassword",
			Status:         UserStatusActive,
			CreateTime:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		if diff := cmp.Diff(expected, out); diff != "" {
			t.Fatalf("unexpected diff: %v", diff)
		}
	})
}

func TestUserWriter_Update(t *testing.T) {
	t.Run("when the user is valid", func(t *testing.T) {
		ctx := context.Background()

		new := &User{
			ID:             uuid.MustParse("18cdf53a-4c02-46df-869a-ba12f2e90d35"),
			OrganizationID: uuid.MustParse("99f85dd9-6df3-4320-8acd-d665a459be38"),
			DisplayName:    "test",
			Email:          "test@test.com",
			HashedPassword: "hashedpassword",
			Status:         UserStatusActive,
			CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		m := &userStorageMock.Repository{
			UpdateFunc: func(ctx context.Context, u *storage.User, fields []storage.UserField) (*storage.User, error) {
				return u, nil
			},
		}

		w := NewUserWriter(func() time.Time { return time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC) }, m)
		out, err := w.Update(ctx, new, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := &User{
			ID:             uuid.MustParse("18cdf53a-4c02-46df-869a-ba12f2e90d35"),
			OrganizationID: uuid.MustParse("99f85dd9-6df3-4320-8acd-d665a459be38"),
			DisplayName:    "test",
			Email:          "test@test.com",
			HashedPassword: "hashedpassword",
			Status:         UserStatusActive,
			CreateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		if diff := cmp.Diff(expected, out); diff != "" {
			t.Fatalf("unexpected diff: %v", diff)
		}
	})
}

func TestUserWriter_Delete(t *testing.T) {
	type fields struct {
		writer storage.UserWriter
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
			w := &UserWriter{
				writer: tt.fields.writer,
			}
			if err := w.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("UserWriter.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
