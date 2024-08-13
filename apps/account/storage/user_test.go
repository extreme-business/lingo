package storage_test

import (
	"errors"
	"testing"

	"github.com/extreme-business/lingo/apps/account/storage"
	"github.com/google/go-cmp/cmp"
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
