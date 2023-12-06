// Package clock provides a simple clock type that can be used to get the current time.
// Its main purpose is to make it easier to test code that depends on the current time.
package clock

import "time"

// Now returns the current time.
type Now func() time.Time

// New returns a new clock that uses the given location and Now function.
func New(loc *time.Location, f Now) Now {
	return func() time.Time { return f().In(loc) }
}

// Default returns a new clock that uses time.Now and time.UTC.
func Default() Now { return New(time.UTC, time.Now) }
