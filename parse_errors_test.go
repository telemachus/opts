package opts_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/opts"
)

func TestParseAlreadyParsedError(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		firstCall  func(*opts.Group, []string) error
		secondCall func(*opts.Group, []string) error
	}{
		"Parse then Parse": {
			firstCall:  (*opts.Group).Parse,
			secondCall: (*opts.Group).Parse,
		},
		"Parse then ParseKnown": {
			firstCall: (*opts.Group).Parse,
			secondCall: func(g *opts.Group, args []string) error {
				_, err := g.ParseKnown(args)
				return err
			},
		},
		"ParseKnown then Parse": {
			firstCall: func(g *opts.Group, args []string) error {
				_, err := g.ParseKnown(args)
				return err
			},
			secondCall: (*opts.Group).Parse,
		},
		"ParseKnown then ParseKnown": {
			firstCall: func(g *opts.Group, args []string) error {
				_, err := g.ParseKnown(args)
				return err
			},
			secondCall: func(g *opts.Group, args []string) error {
				_, err := g.ParseKnown(args)
				return err
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			og := opts.NewGroup("test")
			var b bool
			og.Bool(&b, "verbose")

			// First call should succeed.
			err := tc.firstCall(og, []string{"-verbose"})
			if err != nil {
				t.Fatalf("first call failed: %v", err)
			}

			// Second call should return ErrAlreadyParsed.
			err = tc.secondCall(og, []string{"-verbose"})
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

	err := og.Parse([]string{"-verbose", "extra", "args"})

	var uae *opts.UnexpectedArgsError
	if !errors.As(err, &uae) {
		t.Errorf("expected UnexpectedArgsError, got %T: %v", err, err)
		return
	}

	expectedArgs := []string{"extra", "args"}
	if len(uae.Args) != len(expectedArgs) {
		t.Errorf("expected %d unexpected args, got %d", len(expectedArgs), len(uae.Args))
		return
	}

	for i, expected := range expectedArgs {
		if uae.Args[i] != expected {
			t.Errorf("expected arg[%d] = %q, got %q", i, expected, uae.Args[i])
		}
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
			err := og.Parse(tc.args)

			var mve *opts.MissingValueError
			if !errors.As(err, &mve) {
				t.Errorf("want MissingValueError, got %T: %v", err, err)
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
			err := og.Parse(tc.args)

			var ive *opts.InvalidValueError
			if !errors.As(err, &ive) {
				t.Errorf("expected InvalidValueError, got %T: %v", err, err)
				return
			}

			if ive.Err == nil {
				t.Error("InvalidValueError.Err should not be nil")
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
			err := og.Parse(tc.args)
			if err == nil {
				t.Error("expected an error but got none")
			}
		})
	}
}

func TestParseKnownReturnsExtraArgs(t *testing.T) {
	t.Parallel()

	og := opts.NewGroup("test")
	var verbose bool
	og.Bool(&verbose, "verbose")

	args := []string{"-verbose", "extra", "args"}
	remaining, err := og.ParseKnown(args)
	if err != nil {
		t.Errorf("ParseKnown should not error on extra args, got: %v", err)
	}

	if !verbose {
		t.Error("verbose should be true")
	}

	if diff := cmp.Diff(args[1:], remaining); diff != "" {
		t.Errorf("ParseKnown remaining args (-want +got):\n%s", diff)
	}
}

func TestInvalidValueErrorWrapping(t *testing.T) {
	t.Parallel()

	og := opts.NewGroup("test")
	var i int
	og.Int(&i, "count", 0)

	err := og.Parse([]string{"-count", "not-a-number"})
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
