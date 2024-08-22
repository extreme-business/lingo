package validate_test

import (
	"errors"
	"testing"

	"github.com/extreme-business/lingo/pkg/validate"
)

func TestBytesMinLength(t *testing.T) {
	t.Run("should return nil if the byte slice is at least n characters long", func(t *testing.T) {
		f := validate.BytesMinLength("field", 5)
		if err := f([]byte("12345")); err != nil {
			t.Fatalf("BytesMinLength() error = %v", err)
		}
	})

	t.Run("should return ErrBytesMinLength if the byte slice is less than n characters long", func(t *testing.T) {
		f := validate.BytesMinLength("field", 5)
		if err := f([]byte("1234")); !errors.Is(err, validate.ErrBytesMinLength) {
			t.Fatalf("BytesMinLength() error = %v", err)
		}
	})
}

func TestByteMaxLength(t *testing.T) {
	t.Run("should return nil if the byte slice is at most n characters long", func(t *testing.T) {
		f := validate.ByteMaxLength("field", 5)
		if err := f([]byte("12345")); err != nil {
			t.Fatalf("ByteMaxLength() error = %v", err)
		}
	})

	t.Run("should return ErrBytesMaxLength if the byte slice is more than n characters long", func(t *testing.T) {
		f := validate.ByteMaxLength("field", 5)
		if err := f([]byte("123456")); !errors.Is(err, validate.ErrBytesMaxLength) {
			t.Fatalf("ByteMaxLength() error = %v", err)
		}
	})
}
