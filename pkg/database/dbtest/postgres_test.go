package dbtest_test

import (
	"testing"

	"github.com/extreme-business/lingo/pkg/database/dbtest"
)

func TestSanitizeDBName(t *testing.T) {
	t.Run("should return a sanitized name", func(t *testing.T) {
		got := dbtest.SanitizeDBName("test")
		want := "test"
		if got != want {
			t.Fatalf("SanitizeDBName() = %v, want %v", got, want)
		}
	})

	t.Run("should return a sanitized name with invalid characters removed", func(t *testing.T) {
		type Test struct {
			name string
			want string
		}
		tests := []Test{
			{"test@123", "test_123"},
			{"_-_test@123@456@789@012", "_-_test_123_456_789_012"},
		}
		for _, test := range tests {
			got := dbtest.SanitizeDBName(test.name)
			if got != test.want {
				t.Fatalf("SanitizeDBName() = %v, want %v", got, test.want)
			}
		}
	})
}
