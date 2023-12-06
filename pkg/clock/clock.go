// Package clock provides a simple clock type that can be used to get the current time.
// Its main purpose is to make it easier to test code that depends on the current time.
package clock

import "time"

// Clock can be used to get the current time.
type Clock struct {
	Loc     *time.Location
	NowFunc func() time.Time
}

func (c *Clock) Now() time.Time { return c.NowFunc().In(c.Loc) }

func New(loc *time.Location) *Clock {
	return &Clock{
		Loc:     loc,
		NowFunc: time.Now,
	}
}
