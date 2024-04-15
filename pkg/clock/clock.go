// Package clock provides a simple clock type that can be used to get the current time.
// Its main purpose is to make it easier to test code that depends on the current time.
package clock

import "time"

// Clock can be used to get the current time.
type Clock struct {
	Loc     *time.Location
	nowFunc func() time.Time
}

// Now returns the current time in the clock's location.
func (c *Clock) Now() time.Time { return c.nowFunc().In(c.Loc) }

func New(loc *time.Location, nowFunc func() time.Time) *Clock {
	return &Clock{
		Loc:     loc,
		nowFunc: nowFunc,
	}
}

// Default returns a new clock that uses time.Now and time.UTC.
func Default() *Clock { return New(time.UTC, time.Now) }
