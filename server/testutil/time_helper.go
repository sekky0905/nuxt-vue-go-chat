package testutil

import (
	"time"
)

// reference
// https://medium.com/@timakin/go-api-testing-173b97fb23ec

// FakeTime is mock of time.
var FakeTime time.Time

// SetFakeTime sets FakeTime.
func SetFakeTime(t time.Time) {
	FakeTime = t
}

// ResetFakeTime resets FakeTime.
func ResetFakeTime() {
	FakeTime = time.Time{}
}

// TimeNow returns time now.
func TimeNow() time.Time {
	if !FakeTime.IsZero() {
		return FakeTime
	}
	return time.Now()
}
