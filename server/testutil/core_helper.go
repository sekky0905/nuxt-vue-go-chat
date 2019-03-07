package testutil

import "testing"

// Errorf notify error.
func Errorf(tb testing.TB, want, got interface{}) {
	tb.Helper()
	tb.Errorf("want = %#v, got = %#v", want, got)
}
