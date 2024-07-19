package storage_test

import (
	"errors"
	"testing"

	"github.com/extreme-business/lingo/apps/account/storage"
	"github.com/google/go-cmp/cmp"
)

func TestOrganizationFields(t *testing.T) {
	t.Run("should return the fields", func(t *testing.T) {
		got := storage.OrganizationFields()
		want := []storage.OrganizationField{
			storage.OrganizationID,
			storage.OrganizationLegalName,
			storage.OrganizationSlug,
			storage.OrganizationCreateTime,
			storage.OrganizationUpdateTime,
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("OrganizationFields() mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestOrganizationOrderBy_Validate(t *testing.T) {
	tests := []struct {
		name string
		o    storage.OrganizationOrderBy
		err  error
	}{
		{
			name: "empty",
			o:    storage.OrganizationOrderBy{},
			err:  nil,
		},
		{
			name: "unknown field",
			o:    storage.OrganizationOrderBy{{Field: "invalid"}},
			err:  storage.ErrUnknownOrganizationField,
		},
		{
			name: "empty field",
			o:    storage.OrganizationOrderBy{{Field: ""}},
			err:  storage.ErrEmptyOrganizationSortField,
		},
		{
			name: "valid field and ascending direction",
			o:    storage.OrganizationOrderBy{{Field: storage.OrganizationID, Direction: storage.ASC}},
			err:  nil,
		},
		{
			name: "valid field and descending direction",
			o:    storage.OrganizationOrderBy{{Field: storage.OrganizationID, Direction: storage.DESC}},
			err:  nil,
		},
		{
			name: "invalid direction",
			o:    storage.OrganizationOrderBy{{Field: storage.OrganizationID, Direction: "invalid"}},
			err:  storage.ErrInvalidOrganizationSortDirection,
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
