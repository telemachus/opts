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

			err := og.Parse(tc.args)
			if err != nil {
				t.Fatalf("og.Parse(%v) returns err == %v; want nil", tc.args, err)
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

			err := og.Parse(tc.args)
			if err == nil {
				t.Fatalf("og.Parse(%v) returns err == nil; want UnexpectedArgsError", tc.args)
			}

			var uae *opts.UnexpectedArgsError
			if !errors.As(err, &uae) {
				t.Fatalf("og.Parse(%v) returns err = %v; want UnexpectedArgsError", tc.args, err)
			}

			if len(uae.Args) == 0 {
				t.Errorf("og.Parse(%v) returns empty args; want non-empty slice", tc.args)
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
				t.Fatalf("og.ParseKnown(%v) returns err == %v; want nil", tc.args, err)
			}

			if diff := cmp.Diff(tc.postArgs, remaining); diff != "" {
				t.Errorf("og.ParseKnown(%v); (-want +got):\n%s", tc.args, diff)
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

	err := og.Parse(args)
	if err != nil {
		t.Errorf("og.Parse(%v) returns err = %v; want nil", args, err)
	}

	if cfg.verbose != true {
		t.Errorf("og.Parse(%v) assigns %v to cfg.verbose; want true", args, cfg.verbose)
	}
	if cfg.name != args[2] {
		t.Errorf("og.Parse(%v) assigns %v to cfg.name; want %q", args, cfg.name, args[2])
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
		t.Errorf("og.ParseKnown(%v) returns err = %v; want nil", args, err)
	}

	expectedRemaining := []string{"extra", "args"}
	if diff := cmp.Diff(expectedRemaining, remaining); diff != "" {
		t.Errorf("og.ParseKnown(%v); (-want +got):\n%s", args, diff)
	}

	if cfg.verbose != true {
		t.Errorf("og.ParseKnown(%v) assigns %t to cfg.verbose; want true", args, cfg.verbose)
	}
	if cfg.name != "foobar" {
		t.Errorf("og.ParseKnown(%v) assigns %q to cfg.name; want %q", args, cfg.name, "foobar")
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

			err := og.Parse(tc.args)
			if err == nil {
				t.Errorf("og.Parse(%v) returns err == nil; want error", tc.args)
			}

			og2 := opts.NewGroup("test-parsing")
			_, err2 := og2.ParseKnown(tc.args)
			if err2 == nil {
				t.Errorf("og.ParseKnown(%v) returns err == nil; want error", tc.args)
			}
		})
	}
}

func TestParseRetryAfterFailure(t *testing.T) {
	t.Parallel()

	og := opts.NewGroup("test")
	var s string
	defValue := "default"
	og.String(&s, "name", defValue)

	badArgs := []string{"--unknown", "value"}
	goodArgs := []string{"--name", "success"}

	// First attempt should fail.
	err := og.Parse(badArgs)
	if err == nil {
		t.Fatalf("og.Parse(%v) returns err == nil; want error", badArgs)
	}

	if s != defValue {
		t.Errorf("og.Parse(%v) leaves s as %q; want %q", badArgs, s, defValue)
	}

	// Second attempt should succeed.
	err = og.Parse(goodArgs)
	if err != nil {
		t.Fatalf("og.Parse(%v) returns err == %v; want nil", goodArgs, err)
	}

	if s != goodArgs[1] {
		t.Errorf("og.Parse(%v) leaves s as %q; want %q", goodArgs, s, goodArgs[1])
	}
}

func TestParseKnownRetryAfterFailure(t *testing.T) {
	t.Parallel()

	og := opts.NewGroup("test")
	var s string
	defValue := "default"
	og.String(&s, "name", defValue)

	badArgs := []string{"--unknown", "value"}
	goodArgs := []string{"--name", "success", "arg1", "arg2"}

	// First attempt should fail.
	_, err := og.ParseKnown(badArgs)
	if err == nil {
		t.Fatalf("og.ParseKnown(%v) returns err == nil; want error", badArgs)
	}

	if s != defValue {
		t.Errorf("og.ParseKnown(%v) leaves s as %q; want %q", badArgs, s, defValue)
	}

	// Second attempt should succeed.
	remaining, err := og.ParseKnown(goodArgs)
	if err != nil {
		t.Fatalf("og.ParseKnown(%v) returns err == %v; want nil", goodArgs, err)
	}

	if s != goodArgs[1] {
		t.Errorf("og.ParseKnown(%v) leaves s as %q; want %q", goodArgs, s, goodArgs[1])
	}

	expectedRemaining := []string{"arg1", "arg2"}
	if diff := cmp.Diff(expectedRemaining, remaining); diff != "" {
		t.Errorf("og.ParseKnown(%v); (-want +got):\n%s", goodArgs, diff)
	}
}

func TestAlreadyParsedError(t *testing.T) {
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

			args := []string{"-verbose"}
			og := opts.NewGroup("test")
			var b bool
			og.Bool(&b, "verbose")

			// First call should succeed.
			err := tc.firstCall(og, args)
			if err != nil {
				t.Errorf("Parse(%v) returns %v; want no error", args, err)
			}

			// Second call should return ErrAlreadyParsed.
			err = tc.secondCall(og, args)
			if !errors.Is(err, opts.ErrAlreadyParsed) {
				t.Errorf("Parse(%v) returns %v; want ErrAlreadyParsed", args, err)
			}
		})
	}
}

func TestParseUnexpectedArgsError(t *testing.T) {
	t.Parallel()

	args := []string{"-verbose", "extra", "args"}
	og := opts.NewGroup("test")
	var b bool
	og.Bool(&b, "verbose")

	err := og.Parse(args)

	var uae *opts.UnexpectedArgsError
	if !errors.As(err, &uae) {
		t.Fatalf("Parse(%v) returns %T; want UnexpectedArgsError", args, err)
	}

	expectedArgs := []string{"extra", "args"}
	if diff := cmp.Diff(uae.Args, expectedArgs); diff != "" {
		t.Errorf("og.ParseKnown(%v); (-want +got):\n%s", args, diff)
	}
}
