package opts_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/telemachus/opts"
)

func TestParseAlreadyParsedError(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		firstCall  func(*opts.Group, []string) ([]string, error)
		secondCall func(*opts.Group, []string) ([]string, error)
	}{
		"Parse then Parse": {
			firstCall:  (*opts.Group).Parse,
			secondCall: (*opts.Group).Parse,
		},
		"ParseKnown then ParseKnown": {
			firstCall:  (*opts.Group).ParseKnown,
			secondCall: (*opts.Group).ParseKnown,
		},
		"Parse then ParseKnown": {
			firstCall:  (*opts.Group).Parse,
			secondCall: (*opts.Group).ParseKnown,
		},
		"ParseKnown then Parse": {
			firstCall:  (*opts.Group).ParseKnown,
			secondCall: (*opts.Group).Parse,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			og := opts.NewGroup("test")
			var b bool
			og.Bool(&b, "verbose")

			// First call should succeed
			_, err := tc.firstCall(og, []string{"-verbose"})
			if err != nil {
				t.Fatalf("first call failed: %v", err)
			}

			// Second call should return ErrAlreadyParsed
			_, err = tc.secondCall(og, []string{"-verbose"})
			if !errors.Is(err, opts.ErrAlreadyParsed) {
				t.Errorf("expected ErrAlreadyParsed, got %v", err)
			}
		})
	}
}

func TestParseUnexpectedArgsError(t *testing.T) {
	t.Parallel()

	og := opts.NewGroup("test")
	var b bool
	og.Bool(&b, "verbose")

	_, err := og.Parse([]string{"-verbose", "extra", "args"})
	if !errors.Is(err, opts.ErrUnexpectedArgs) {
		t.Errorf("expected ErrUnexpectedArgs, got %v", err)
	}
}

func TestParseMissingValueError(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		setup func() *opts.Group
		args  []string
	}{
		"string option missing value": {
			setup: func() *opts.Group {
				og := opts.NewGroup("test")
				var s string
				og.String(&s, "file", "default")
				return og
			},
			args: []string{"-file"},
		},
		"int option missing value": {
			setup: func() *opts.Group {
				og := opts.NewGroup("test")
				var i int
				og.Int(&i, "count", 0)
				return og
			},
			args: []string{"--count"},
		},
		"equals with empty value": {
			setup: func() *opts.Group {
				og := opts.NewGroup("test")
				var s string
				og.String(&s, "file", "default")
				return og
			},
			args: []string{"--file="},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			og := tc.setup()
			_, err := og.Parse(tc.args)
			if !errors.Is(err, opts.ErrMissingValue) {
				t.Errorf("expected ErrMissingValue, got %v", err)
			}
		})
	}
}

func TestParseInvalidValueError(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		setup func() *opts.Group
		args  []string
	}{
		"int conversion error": {
			setup: func() *opts.Group {
				og := opts.NewGroup("test")
				var i int
				og.Int(&i, "count", 0)
				return og
			},
			args: []string{"-count", "not-a-number"},
		},
		"float conversion error": {
			setup: func() *opts.Group {
				og := opts.NewGroup("test")
				var f float64
				og.Float64(&f, "value", 0.0)
				return og
			},
			args: []string{"--value=invalid-float"},
		},
		"uint negative value": {
			setup: func() *opts.Group {
				og := opts.NewGroup("test")
				var u uint
				og.Uint(&u, "count", 0)
				return og
			},
			args: []string{"-count", "-5"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			og := tc.setup()
			_, err := og.Parse(tc.args)

			var ive *opts.InvalidValueError
			if !errors.As(err, &ive) {
				t.Errorf("expected InvalidValueError, got %T: %v", err, err)
				return
			}

			if ive.Err == nil {
				t.Error("InvalidValueError.Err should not be nil")
			}

			// Check that it's a reasonable underlying error type
			if !errors.Is(ive.Err, strconv.ErrSyntax) && !errors.Is(ive.Err, strconv.ErrRange) {
				t.Logf("Underlying error: %T: %v", ive.Err, ive.Err)
			}
		})
	}
}

func TestParseGeneralErrors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		setup func() *opts.Group
		args  []string
	}{
		"unknown option": {
			setup: func() *opts.Group {
				og := opts.NewGroup("test")
				var b bool
				og.Bool(&b, "verbose")
				return og
			},
			args: []string{"-unknown"},
		},
		"bool with equals": {
			setup: func() *opts.Group {
				og := opts.NewGroup("test")
				var b bool
				og.Bool(&b, "verbose")
				return og
			},
			args: []string{"--verbose=true"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			og := tc.setup()
			_, err := og.Parse(tc.args)
			if err == nil {
				t.Error("expected an error but got none")
			}
		})
	}
}

func TestParseKnownDoesNotErrorOnExtraArgs(t *testing.T) {
	t.Parallel()

	og := opts.NewGroup("test")
	var verbose bool
	og.Bool(&verbose, "verbose")

	remaining, err := og.ParseKnown([]string{"-verbose", "extra", "args"})
	if err != nil {
		t.Errorf("ParseKnown should not error on extra args, got: %v", err)
	}

	if !verbose {
		t.Error("verbose should be true")
	}

	expectedRemaining := []string{"extra", "args"}
	if len(remaining) != len(expectedRemaining) {
		t.Errorf("expected %d remaining args, got %d", len(expectedRemaining), len(remaining))
	}

	for i, expected := range expectedRemaining {
		if i >= len(remaining) || remaining[i] != expected {
			t.Errorf("expected remaining[%d] = %q, got %q", i, expected, remaining[i])
		}
	}
}

func TestInvalidValueErrorWrapping(t *testing.T) {
	t.Parallel()

	og := opts.NewGroup("test")
	var i int
	og.Int(&i, "count", 0)

	_, err := og.Parse([]string{"-count", "not-a-number"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var ive *opts.InvalidValueError
	if !errors.As(err, &ive) {
		t.Fatalf("expected InvalidValueError, got %T", err)
	}

	// Test that Unwrap works
	unwrapped := errors.Unwrap(ive)
	if unwrapped == nil {
		t.Error("Unwrap should return the underlying error")
	}

	if !errors.Is(unwrapped, ive.Err) {
		t.Error("Unwrap should return the same error as Err field")
	}
}
