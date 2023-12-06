package uuidgen_test

import (
	"testing"

	"github.com/dwethmar/lingo/pkg/uuidgen"
)

func TestDefault(t *testing.T) {
	t.Run("create default generator", func(t *testing.T) {
		g := uuidgen.Default()
		if g == nil {
			t.Errorf("Default() = %v, want non-nil", g)
		}

		if g().String() == "" {
			t.Errorf("Default() = %v, want non-empty", g().String())
		}
	})
}
