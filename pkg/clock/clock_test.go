package clock_test

import (
	"testing"
	"time"

	"github.com/dwethmar/lingo/pkg/clock"
)

func TestClock_Now(t *testing.T) {
	t.Run("get current time", func(t *testing.T) {
		c := clock.New(time.UTC, time.Now)
		now := c.Now()
		if now.IsZero() {
			t.Errorf("Clock.Now() = %v, want non zero", now)
		}
	})

	t.Run("get current time in a different location", func(t *testing.T) {
		loc, _ := time.LoadLocation("America/New_York")
		c := clock.New(loc, time.Now)
		now := c.Now()
		if now.IsZero() {
			t.Errorf("Clock.Now() = %v, want non zero", now)
		}

		if now.Location() != loc {
			t.Errorf("Clock.Now() location = %v, want %v", now.Location(), loc)
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("create new clock", func(t *testing.T) {
		c := clock.New(time.UTC, time.Now)
		if c == nil {
			t.Errorf("New() = %v, want non-nil", c)
		}
	})
}

func TestDefault(t *testing.T) {
	t.Run("create default clock", func(t *testing.T) {
		c := clock.Default()
		if c == nil {
			t.Errorf("Default() = %v, want non-nil", c)
		}

		// Check if the location is set to UTC
		if l := c.Location(); l != time.UTC {
			t.Errorf("Default() location = %v, want %v", l, time.UTC)
		}
	})
}
