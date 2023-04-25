package timesource

import (
	"time"
)

// TimeSource provides current time.
type TimeSource interface {
	Now() time.Time
}

// DefaultTimeSource serves as a UTC time source.
type DefaultTimeSource struct{}

// Now returns current UTC time.
func (d DefaultTimeSource) Now() time.Time {
	return time.Now().UTC()
}
