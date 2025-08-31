package opts_test

import (
	"errors"
	"testing"
)

func checkErrorAs[T error](t *testing.T, err error) {
	t.Helper()

	var target T
	if !errors.As(err, &target) {
		// Using %T on a nil pointer of a specific type gives the type name.
		t.Errorf("expected error of type %T, but got %T: %v", target, err, err)
	}
}
