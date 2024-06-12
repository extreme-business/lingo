package validate_test

import (
	"errors"
	"testing"

	"github.com/extreme-business/lingo/pkg/validate"
	"github.com/google/uuid"
)

func TestUUIDValidator_Validate(t *testing.T) {
	type args struct {
		s uuid.UUID
	}
	tests := []struct {
		name    string
		v       validate.UUIDValidator
		args    args
		wantErr bool
	}{
		{
			name:    "should return an error if the uuid is nil",
			v:       validate.UUIDValidator{validate.UUIDIsNotNil("test")},
			args:    args{s: uuid.Nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.Validate(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("UUIDValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUUIDIsNotNil(t *testing.T) {
	t.Run("should return an error if the uuid is nil", func(t *testing.T) {
		v := validate.UUIDIsNotNil("test")
		if err := v(uuid.Nil); err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("should return no error if the uuid is not nil", func(t *testing.T) {
		v := validate.UUIDIsNotNil("test")
		if err := v(uuid.MustParse("7fb3d880-1db0-464e-b062-a9896cb9bf6c")); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("error matches field, message and validate.ErrUUIDIsNil", func(t *testing.T) {
		v := validate.UUIDIsNotNil("test")
		err := v(uuid.Nil)
		if err == nil {
			t.Error("expected an error")
		}

		if err.Field() != "test" {
			t.Errorf("expected field to be test, got %s", err.Field())
		}

		if err.Error() != "test: should not be nil" {
			t.Errorf("expected error to be test: should not be nil, got %s", err.Error())
		}

		if !errors.Is(err, validate.ErrUUIDIsNil) {
			t.Errorf("expected wrapped error to be ErrUUIDIsNil, got %v", err.Unwrap())
		}
	})
}
