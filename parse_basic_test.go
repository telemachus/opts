package opts_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/opts"
)

func TestParseStrictNoOptions(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args []string
	}{
		"Empty args": {args: []string{}},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			og := opts.NewGroup("test-parsing")

			remaining, err := og.Parse(tc.args)
			if err != nil {
				t.Fatalf("after remaining, err := og.Parse(%v), err == %v; want nil", tc.args, err)
			}

			if len(remaining) != 0 {
				t.Errorf("after og.Parse(%v), remaining = %v; want empty slice", tc.args, remaining)
			}
		})
	}
}

func TestParseStrictWithRemainingArgs(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args []string
	}{
		"-- followed by args": {
			args: []string{"--", "foo", "bar"},
		},
		"- followed by args": {
			args: []string{"-", "foo", "bar"},
		},
		"non-option args": {
			args: []string{"foo", "bar"},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			og := opts.NewGroup("test-parsing")

			remaining, err := og.Parse(tc.args)
			if err == nil {
				t.Fatalf("after remaining, err := og.Parse(%v), err == nil; want ErrUnexpectedArgs", tc.args)
			}

			if !errors.Is(err, opts.ErrUnexpectedArgs) {
				t.Errorf("after og.Parse(%v), err = %v; want ErrUnexpectedArgs", tc.args, err)
			}

			if len(remaining) == 0 {
				t.Errorf("after og.Parse(%v), remaining = %v; want non-empty slice", tc.args, remaining)
			}
		})
	}
}

func TestParseKnownNoOptions(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args     []string
		postArgs []string
	}{
		"Empty args": {args: []string{}, postArgs: []string{}},
		"-- should not be in remaining": {
			args:     []string{"--", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
		},
		"- should be in remaining": {
			args:     []string{"-", "foo", "bar"},
			postArgs: []string{"-", "foo", "bar"},
		},
		"non-option args": {
			args:     []string{"foo", "bar"},
			postArgs: []string{"foo", "bar"},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			og := opts.NewGroup("test-parsing")

			remaining, err := og.ParseKnown(tc.args)
			if err != nil {
				t.Fatalf("after remaining, err := og.ParseKnown(%v), err == %v; want nil", tc.args, err)
			}

			if diff := cmp.Diff(tc.postArgs, remaining); diff != "" {
				t.Errorf("after og.ParseKnown(%v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

func TestParseStrictWithOptions(t *testing.T) {
	t.Parallel()

	args := []string{"-V", "--name", "foobar"}
	cfg := struct {
		name    string
		verbose bool
	}{}

	og := opts.NewGroup("test-parsing")
	og.Bool(&cfg.verbose, "V")
	og.String(&cfg.name, "name", "default")

	remaining, err := og.Parse(args)
	if err != nil {
		t.Errorf("after remaining, err := og.Parse(%v), err = %v; want nil", args, err)
	}

	if len(remaining) != 0 {
		t.Errorf("after og.Parse(%v), remaining = %v; want empty slice", args, remaining)
	}

	if cfg.verbose != true {
		t.Errorf("after og.Parse(%v), cfg.verbose = %v; want true", args, cfg.verbose)
	}
	if cfg.name != args[2] {
		t.Errorf("after og.Parse(%v), cfg.name = %v; want %q", args, cfg.name, args[2])
	}
}

func TestParseKnownWithOptions(t *testing.T) {
	t.Parallel()

	args := []string{"-V", "--name", "foobar", "extra", "args"}
	cfg := struct {
		name    string
		verbose bool
	}{}

	og := opts.NewGroup("test-parsing")
	og.Bool(&cfg.verbose, "V")
	og.String(&cfg.name, "name", "default")

	remaining, err := og.ParseKnown(args)
	if err != nil {
		t.Errorf("after remaining, err := og.ParseKnown(%v), err = %v; want nil", args, err)
	}

	expectedRemaining := []string{"extra", "args"}
	if diff := cmp.Diff(expectedRemaining, remaining); diff != "" {
		t.Errorf("after og.ParseKnown(%v); (-want +got):\n%s", args, diff)
	}

	if cfg.verbose != true {
		t.Errorf("after og.ParseKnown(%v), cfg.verbose = %v; want true", args, cfg.verbose)
	}
	if cfg.name != "foobar" {
		t.Errorf("after og.ParseKnown(%v), cfg.name = %v; want %q", args, cfg.name, "foobar")
	}
}

func TestParseUndefinedOptions(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args []string
	}{
		"Unknown single-dash option": {args: []string{"-x"}},
		"Unknown double-dash option": {args: []string{"--unknown"}},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			og := opts.NewGroup("test-parsing")

			_, err := og.Parse(tc.args)
			if err == nil {
				t.Errorf("after og.Parse(%v), err == nil; want error", tc.args)
			}

			og2 := opts.NewGroup("test-parsing")
			_, err2 := og2.ParseKnown(tc.args)
			if err2 == nil {
				t.Errorf("after og.ParseKnown(%v), err == nil; want error", tc.args)
			}
		})
	}
}

func TestReparseAfterFailure(t *testing.T) {
	t.Parallel()

	t.Run("Parse strict", func(t *testing.T) {
		testParseRetryAfterFailure(t, (*opts.Group).Parse, []string{"--name", "success"})
	})

	t.Run("ParseKnown", func(t *testing.T) {
		testParseRetryAfterFailure(t, (*opts.Group).ParseKnown, []string{"--name", "success", "arg1", "arg2"})
	})
}

func testParseRetryAfterFailure(t *testing.T, parseFunc func(*opts.Group, []string) ([]string, error), goodArgs []string) {
	t.Helper()
	t.Parallel()

	og := opts.NewGroup("test")
	var s string
	defValue := "default"
	og.String(&s, "name", defValue)

	badArgs := []string{"--unknown", "value"}

	// First attempt should fail.
	_, err := parseFunc(og, badArgs)
	if err == nil {
		t.Fatalf("after parseFunc(%v), err == nil; want error", badArgs)
	}

	if s != defValue {
		t.Errorf("after failed parse, s == %q; want %q", s, defValue)
	}

	// Second attempt should succeed.
	remaining, err := parseFunc(og, goodArgs)
	if err != nil {
		t.Fatalf("after parseFunc(%v), err == %v; want nil", goodArgs, err)
	}

	if s != goodArgs[1] {
		t.Errorf("after parseFunc(%v), s == %q; want %q", goodArgs, s, goodArgs[1])
	}

	// For ParseKnown, check remaining args; for Parse, should be empty.
	expectedRemaining := goodArgs[2:]
	if diff := cmp.Diff(expectedRemaining, remaining); diff != "" {
		t.Errorf("after parseFunc(%v); (-want +got):\n%s", goodArgs, diff)
	}
}
