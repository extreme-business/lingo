package storage_test

import (
	"errors"
	"testing"
	"time"

	"github.com/extreme-business/lingo/cmd/account/domain"
	"github.com/extreme-business/lingo/cmd/account/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestUserFields(t *testing.T) {
	t.Run("get user fields", func(t *testing.T) {
		got := storage.UserFields()
		want := []storage.UserField{
			storage.UserID,
			storage.UserOrganizationID,
			storage.UserDisplayName,
			storage.UserEmail,
			storage.UserPassword,
			storage.UserStatus,
			storage.UserCreateTime,
			storage.UserUpdateTime,
			storage.UserDeleteTime,
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("UserFields() mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestUser_ToDomain(t *testing.T) {
	type fields struct {
		ID             uuid.UUID
		OrganizationID uuid.UUID
		DisplayName    string
		Email          string
		Password       string
		CreateTime     time.Time
		UpdateTime     time.Time
	}
	type args struct {
		in *domain.User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *domain.User
	}{
		{
			name: "convert to domain",
			fields: fields{
				ID:             uuid.MustParse("0d322d31-c960-497b-ada0-d3ffd1bded8f"),
				OrganizationID: uuid.MustParse("95a2122b-3591-4f42-bfd2-c5b8d3f8c30b"),
				DisplayName:    "display name",
				Email:          "email",
				Password:       "password",
				CreateTime:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdateTime:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				in: &domain.User{
					ID:             uuid.MustParse("0d322d31-c960-497b-ada0-d3ffd1bded8f"),
					OrganizationID: uuid.MustParse("95a2122b-3591-4f42-bfd2-c5b8d3f8c30b"),
					DisplayName:    "display name",
					Email:          "email",
					HashedPassword: "password",
					CreateTime:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdateTime:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			want: &domain.User{
				ID:             uuid.MustParse("0d322d31-c960-497b-ada0-d3ffd1bded8f"),
				OrganizationID: uuid.MustParse("95a2122b-3591-4f42-bfd2-c5b8d3f8c30b"),
				DisplayName:    "display name",
				Email:          "email",
				HashedPassword: "password",
				CreateTime:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdateTime:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &storage.User{
				ID:             tt.fields.ID,
				OrganizationID: tt.fields.OrganizationID,
				DisplayName:    tt.fields.DisplayName,
				Email:          tt.fields.Email,
				HashedPassword: tt.fields.Password,
				CreateTime:     tt.fields.CreateTime,
				UpdateTime:     tt.fields.UpdateTime,
			}
			u.ToDomain(tt.args.in)

			if diff := cmp.Diff(tt.args.in, tt.want); diff != "" {
				t.Errorf("User.ToDomain() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUser_FromDomain(t *testing.T) {
	type fields struct {
		ID             uuid.UUID
		OrganizationID uuid.UUID
		DisplayName    string
		Email          string
		Password       string
		CreateTime     time.Time
		UpdateTime     time.Time
	}
	type args struct {
		in *domain.User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "convert from domain",
			fields: fields{
				ID:             uuid.MustParse("0d322d31-c960-497b-ada0-d3ffd1bded8f"),
				OrganizationID: uuid.MustParse("95a2122b-3591-4f42-bfd2-c5b8d3f8c30b"),
				DisplayName:    "display name",
				Email:          "email",
				Password:       "password",
				CreateTime:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdateTime:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				in: &domain.User{
					ID:             uuid.MustParse("0d322d31-c960-497b-ada0-d3ffd1bded8f"),
					OrganizationID: uuid.MustParse("95a2122b-3591-4f42-bfd2-c5b8d3f8c30b"),
					DisplayName:    "display name",
					Email:          "email",
					HashedPassword: "password",
					CreateTime:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdateTime:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			u := &storage.User{
				ID:             tt.fields.ID,
				OrganizationID: tt.fields.OrganizationID,
				DisplayName:    tt.fields.DisplayName,
				Email:          tt.fields.Email,
				HashedPassword: tt.fields.Password,
				CreateTime:     tt.fields.CreateTime,
				UpdateTime:     tt.fields.UpdateTime,
			}
			u.FromDomain(tt.args.in)
		})
	}
}

func TestUserOrderBy_Validate(t *testing.T) {
	tests := []struct {
		name string
		o    storage.UserOrderBy
		err  error
	}{
		{
			name: "empty",
			o:    storage.UserOrderBy{},
			err:  nil,
		},
		{
			name: "unknown field",
			o:    storage.UserOrderBy{{Field: "invalid"}},
			err:  storage.ErrUserUnknownField,
		},
		{
			name: "empty field",
			o:    storage.UserOrderBy{{Field: ""}},
			err:  storage.ErrEmptyUserSortField,
		},
		{
			name: "valid field and ascending direction",
			o:    storage.UserOrderBy{{Field: storage.UserID, Direction: storage.ASC}},
			err:  nil,
		},
		{
			name: "valid field and descending direction",
			o:    storage.UserOrderBy{{Field: storage.UserID, Direction: storage.DESC}},
			err:  nil,
		},
		{
			name: "invalid direction",
			o:    storage.UserOrderBy{{Field: storage.UserID, Direction: "invalid"}},
			err:  storage.ErrInvalidUserSortDirection,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.o.Validate(); !errors.Is(err, tt.err) {
				t.Errorf("UserOrderBy.Validate() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
