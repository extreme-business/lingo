package uuidgen

import (
	"testing"

	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	t.Run("should return a new generator", func(t *testing.T) {
		g := New(func() uuid.UUID {
			return uuid.New()
		})
		if g == nil {
			t.Error("expected a new generator")
			return
		}

		if g.genFunc == nil {
			t.Error("expected a generator with a genFunc")
			return
		}

		if len(g.genFunc()) == 0 {
			t.Error("expected a generator with a genFunc that returns a uuid")
			return
		}
	})
}

func TestGenerator_New(t *testing.T) {
	t.Run("should return a new uuid", func(t *testing.T) {
		g := Generator{
			genFunc: func() uuid.UUID {
				return uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000000"))
			},
		}

		if g.New() != uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000000")) {
			t.Error("expected a new uuid")
		}
	})
}

func TestDefault(t *testing.T) {
	t.Run("should return a new generator", func(t *testing.T) {
		g := Default()
		if g == nil {
			t.Error("expected a new generator")
			return
		}

		if g.genFunc == nil {
			t.Error("expected a generator with a genFunc")
		}
	})
}
