package opts_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// checkRemaining reports any difference between the leftover arguments that
// og.method(args) returned and the ones that the test wants.
func checkRemaining(t *testing.T, method string, args, want, got []string) {
	t.Helper()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("og.%s(%v); (-want +got):\n%s", method, args, diff)
	}
}
