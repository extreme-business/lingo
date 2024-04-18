package uuidgen_test

import (
	"testing"

	"github.com/dwethmar/lingo/pkg/uuidgen"
	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	t.Run("should return a new generator", func(t *testing.T) {
		g := uuidgen.New(func() uuid.UUID {
			return uuid.New()
		})
		if g == nil {
			t.Error("expected a new generator")
			return
		}
	})
}

func TestGenerator_New(t *testing.T) {
	t.Run("should return a new uuid", func(t *testing.T) {
		if uuidgen.New(func() uuid.UUID { return uuid.Nil }).New() != uuid.Nil {
			t.Error("expected the same uuid")
		}
	})
}

func TestDefault(t *testing.T) {
	t.Run("should return a new generator", func(t *testing.T) {
		g := uuidgen.Default()
		if g == nil {
			t.Error("expected a new generator")
			return
		}
	})
}
