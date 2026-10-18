package opts_test

import (
	"errors"
	"testing"

	"github.com/telemachus/opts"
)

func TestParseAllNoOptions(t *testing.T) {
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

			err := og.ParseAll(tc.args)
			if err != nil {
				t.Fatalf("og.ParseAll(%v) returns err == %v; want nil", tc.args, err)
			}
		})
	}
}

func TestParseAllWithRemainingArgs(t *testing.T) {
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

			// errors.AsType reports false for a nil error, so this also
			// catches a missing error.
			err := og.ParseAll(tc.args)
			uae, ok := errors.AsType[*opts.UnexpectedArgumentsError](err)
			if !ok {
				t.Fatalf("og.ParseAll(%v) returns err = %v; want UnexpectedArgumentsError", tc.args, err)
			}

			if len(uae.Args) == 0 {
				t.Errorf("og.ParseAll(%v) returns empty args; want non-empty slice", tc.args)
			}
		})
	}
}

func TestParseNoOptions(t *testing.T) {
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

			remaining, err := og.Parse(tc.args)
			if err != nil {
				t.Fatalf("og.Parse(%v) returns err == %v; want nil", tc.args, err)
			}

			checkRemaining(t, "Parse", tc.args, tc.postArgs, remaining)
		})
	}
}

func TestParseAllWithOptions(t *testing.T) {
	t.Parallel()

	args := []string{"-V", "--name", "foobar"}
	cfg := struct {
		name    string
		verbose bool
	}{}

	og := opts.NewGroup("test-parsing")
	og.Bool(&cfg.verbose, "V")
	og.String(&cfg.name, "name", "default")

	err := og.ParseAll(args)
	if err != nil {
		t.Errorf("og.ParseAll(%v) returns err = %v; want nil", args, err)
	}

	if !cfg.verbose {
		t.Errorf("og.ParseAll(%v) assigns %v to cfg.verbose; want true", args, cfg.verbose)
	}
	if cfg.name != args[2] {
		t.Errorf("og.ParseAll(%v) assigns %v to cfg.name; want %q", args, cfg.name, args[2])
	}
}

func TestParseWithOptions(t *testing.T) {
	t.Parallel()

	args := []string{"-V", "--name", "foobar", "extra", "args"}
	cfg := struct {
		name    string
		verbose bool
	}{}

	og := opts.NewGroup("test-parsing")
	og.Bool(&cfg.verbose, "V")
	og.String(&cfg.name, "name", "default")

	remaining, err := og.Parse(args)
	if err != nil {
		t.Errorf("og.Parse(%v) returns err = %v; want nil", args, err)
	}

	expectedRemaining := []string{"extra", "args"}
	checkRemaining(t, "Parse", args, expectedRemaining, remaining)

	if !cfg.verbose {
		t.Errorf("og.Parse(%v) assigns %t to cfg.verbose; want true", args, cfg.verbose)
	}
	if cfg.name != "foobar" {
		t.Errorf("og.Parse(%v) assigns %q to cfg.name; want %q", args, cfg.name, "foobar")
	}
}

func TestParseUndefinedOptions(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args []string
	}{
		"Unknown single-dash option":    {args: []string{"-x"}},
		"Unknown double-dash option":    {args: []string{"--unknown"}},
		"No name before equals; single": {args: []string{"-=x"}},
		"No name before equals; double": {args: []string{"--="}},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			og := opts.NewGroup("test-parsing")

			err := og.ParseAll(tc.args)
			if !errors.Is(err, opts.ErrUnknownOption) {
				t.Errorf("og.ParseAll(%v) returns err = %v; want ErrUnknownOption", tc.args, err)
			}

			og2 := opts.NewGroup("test-parsing")
			_, err2 := og2.Parse(tc.args)
			if !errors.Is(err2, opts.ErrUnknownOption) {
				t.Errorf("og.Parse(%v) returns err = %v; want ErrUnknownOption", tc.args, err2)
			}
		})
	}
}

func TestParseAllRetryAfterFailure(t *testing.T) {
	t.Parallel()

	og := opts.NewGroup("test")
	var name string
	defValue := "default"
	og.String(&name, "name", defValue)

	badArgs := []string{"--unknown", "value"}
	goodArgs := []string{"--name", "success"}

	// First attempt should fail.
	err := og.ParseAll(badArgs)
	if err == nil {
		t.Fatalf("og.ParseAll(%v) returns err == nil; want error", badArgs)
	}

	if name != defValue {
		t.Errorf("og.ParseAll(%v) leaves name as %q; want %q", badArgs, name, defValue)
	}

	// Second attempt should succeed.
	err = og.ParseAll(goodArgs)
	if err != nil {
		t.Fatalf("og.ParseAll(%v) returns err == %v; want nil", goodArgs, err)
	}

	if name != goodArgs[1] {
		t.Errorf("og.ParseAll(%v) leaves name as %q; want %q", goodArgs, name, goodArgs[1])
	}
}

func TestParseRetryAfterFailure(t *testing.T) {
	t.Parallel()

	og := opts.NewGroup("test")
	var name string
	defValue := "default"
	og.String(&name, "name", defValue)

	badArgs := []string{"--unknown", "value"}
	goodArgs := []string{"--name", "success", "arg1", "arg2"}

	// First attempt should fail.
	_, err := og.Parse(badArgs)
	if err == nil {
		t.Fatalf("og.Parse(%v) returns err == nil; want error", badArgs)
	}

	if name != defValue {
		t.Errorf("og.Parse(%v) leaves name as %q; want %q", badArgs, name, defValue)
	}

	// Second attempt should succeed.
	remaining, err := og.Parse(goodArgs)
	if err != nil {
		t.Fatalf("og.Parse(%v) returns err == %v; want nil", goodArgs, err)
	}

	if name != goodArgs[1] {
		t.Errorf("og.Parse(%v) leaves name as %q; want %q", goodArgs, name, goodArgs[1])
	}

	expectedRemaining := []string{"arg1", "arg2"}
	checkRemaining(t, "Parse", goodArgs, expectedRemaining, remaining)
}

func TestAlreadyParsedError(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		firstCall  func(*opts.Group, []string) error
		secondCall func(*opts.Group, []string) error
	}{
		"ParseAll then ParseAll": {
			firstCall:  (*opts.Group).ParseAll,
			secondCall: (*opts.Group).ParseAll,
		},
		"ParseAll then Parse": {
			firstCall: (*opts.Group).ParseAll,
			secondCall: func(g *opts.Group, args []string) error {
				_, err := g.Parse(args)
				return err
			},
		},
		"Parse then ParseAll": {
			firstCall: func(g *opts.Group, args []string) error {
				_, err := g.Parse(args)
				return err
			},
			secondCall: (*opts.Group).ParseAll,
		},
		"Parse then Parse": {
			firstCall: func(g *opts.Group, args []string) error {
				_, err := g.Parse(args)
				return err
			},
			secondCall: func(g *opts.Group, args []string) error {
				_, err := g.Parse(args)
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
				t.Errorf("ParseAll(%v) returns %v; want no error", args, err)
			}

			// Second call should return ErrAlreadyParsed.
			err = tc.secondCall(og, args)
			if !errors.Is(err, opts.ErrAlreadyParsed) {
				t.Errorf("ParseAll(%v) returns %v; want ErrAlreadyParsed", args, err)
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

	err := og.ParseAll(args)

	uae, ok := errors.AsType[*opts.UnexpectedArgumentsError](err)
	if !ok {
		t.Fatalf("ParseAll(%v) returns %T; want UnexpectedArgumentsError", args, err)
	}

	expectedArgs := []string{"extra", "args"}
	checkRemaining(t, "ParseAll", args, expectedArgs, uae.Args)
}
