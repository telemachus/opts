package opts_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/telemachus/opts"
)

const (
	defName  = "default"
	defCount = 1
)

type retryCfg struct {
	name    string
	count   int
	verbose bool
}

func newRetryGroup(cfg *retryCfg) *opts.Group {
	og := opts.NewGroup("test-retry")
	og.String(&cfg.name, "name", defName)
	og.String(&cfg.name, "n", defName)
	og.Int(&cfg.count, "count", defCount)
	og.Bool(&cfg.verbose, "verbose")

	return og
}

func parseDiscardingRest(g *opts.Group, args []string) error {
	_, err := g.Parse(args)
	return err
}

var parseMethods = map[string]func(*opts.Group, []string) error{
	"ParseAll": (*opts.Group).ParseAll,
	"Parse":    parseDiscardingRest,
}

func checkDefaults(t *testing.T, cfg *retryCfg, method string, args []string) {
	t.Helper()

	if cfg.name != defName {
		t.Errorf("og.%s(%v) leaves cfg.name as %q; want %q", method, args, cfg.name, defName)
	}
	if cfg.count != defCount {
		t.Errorf("og.%s(%v) leaves cfg.count as %d; want %d", method, args, cfg.count, defCount)
	}
	if cfg.verbose {
		t.Errorf("og.%s(%v) leaves cfg.verbose as true; want false", method, args)
	}
}

// checkRetried checks the values after a successful retry with --verbose:
// verbose is set and the other options keep their defaults.
func checkRetried(t *testing.T, cfg *retryCfg, method string, args []string) {
	t.Helper()

	if !cfg.verbose {
		t.Errorf("og.%s(%v) leaves cfg.verbose as false; want true", method, args)
	}
	if cfg.name != defName {
		t.Errorf("og.%s(%v) leaves cfg.name as %q; want %q", method, args, cfg.name, defName)
	}
	if cfg.count != defCount {
		t.Errorf("og.%s(%v) leaves cfg.count as %d; want %d", method, args, cfg.count, defCount)
	}
}

func TestParseFailureLeavesValuesUnchanged(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		target error
		args   []string
	}{
		"unknown option after valid options": {
			args:   []string{"--name", "bad", "--verbose", "--unknown"},
			target: opts.ErrUnknownOption,
		},
		"invalid value after valid options": {
			args:   []string{"--name", "bad", "--verbose", "--count", "abc"},
			target: strconv.ErrSyntax,
		},
		"missing value after valid options": {
			args:   []string{"--name", "bad", "--verbose", "--count"},
			target: opts.ErrMissingValue,
		},
		"boolean with value after valid options": {
			args:   []string{"--name", "bad", "--count=7", "--verbose=true"},
			target: opts.ErrBooleanWithValue,
		},
		"empty value": {
			args:   []string{"--count="},
			target: opts.ErrMissingValue,
		},
		"empty value after valid options": {
			args:   []string{"--name", "bad", "--verbose", "--count="},
			target: opts.ErrMissingValue,
		},
	}

	for method, parse := range parseMethods {
		for msg, tc := range testCases {
			t.Run(method+"/"+msg, func(t *testing.T) {
				t.Parallel()

				var cfg retryCfg
				og := newRetryGroup(&cfg)

				err := parse(og, tc.args)
				if !errors.Is(err, tc.target) {
					t.Fatalf("og.%s(%v) returns err = %v; want %v", method, tc.args, err, tc.target)
				}

				checkDefaults(t, &cfg, method, tc.args)
			})
		}
	}
}

func TestParseUnexpectedArgsLeavesValuesUnchanged(t *testing.T) {
	t.Parallel()

	args := []string{"--name", "bad", "--count", "7", "--verbose", "extra"}

	var cfg retryCfg
	og := newRetryGroup(&cfg)

	err := og.ParseAll(args)

	if _, ok := errors.AsType[*opts.UnexpectedArgumentsError](err); !ok {
		t.Fatalf("og.ParseAll(%v) returns err = %v; want UnexpectedArgumentsError", args, err)
	}

	checkDefaults(t, &cfg, "ParseAll", args)
}

//revive:disable-next-line:cognitive-complexity
func TestParseRetryDoesNotKeepFailedValues(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		badArgs []string
	}{
		"unknown option after valid options": {
			badArgs: []string{"--name", "bad", "--count", "7", "--unknown"},
		},
		"invalid value after valid options": {
			badArgs: []string{"--name", "bad", "--count", "abc"},
		},
	}

	goodArgs := []string{"--verbose"}

	for method, parse := range parseMethods {
		for msg, tc := range testCases {
			t.Run(method+"/"+msg, func(t *testing.T) {
				t.Parallel()

				var cfg retryCfg
				og := newRetryGroup(&cfg)

				// First attempt should fail.
				if err := parse(og, tc.badArgs); err == nil {
					t.Fatalf("og.%s(%v) returns err == nil; want error", method, tc.badArgs)
				}

				// Second attempt should succeed and keep defaults for options
				// that it does not mention.
				if err := parse(og, goodArgs); err != nil {
					t.Fatalf("og.%s(%v) returns err = %v; want nil", method, goodArgs, err)
				}

				checkRetried(t, &cfg, method, goodArgs)
			})
		}
	}
}

func TestParseRetryAfterUnexpectedArgs(t *testing.T) {
	t.Parallel()

	badArgs := []string{"--name", "bad", "--count", "7", "extra"}
	goodArgs := []string{"--verbose"}

	var cfg retryCfg
	og := newRetryGroup(&cfg)

	// First attempt should fail.
	err := og.ParseAll(badArgs)

	if _, ok := errors.AsType[*opts.UnexpectedArgumentsError](err); !ok {
		t.Fatalf("og.ParseAll(%v) returns err = %v; want UnexpectedArgumentsError", badArgs, err)
	}

	// Second attempt should succeed and keep defaults for options that it
	// does not mention.
	if err := og.ParseAll(goodArgs); err != nil {
		t.Fatalf("og.ParseAll(%v) returns err = %v; want nil", goodArgs, err)
	}

	checkRetried(t, &cfg, "ParseAll", goodArgs)
}

//revive:disable-next-line:cognitive-complexity
func TestParseAliasLastValueWins(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		want string
		args []string
	}{
		"short then long": {
			args: []string{"-n", "first", "--name", "second"},
			want: "second",
		},
		"long then short": {
			args: []string{"--name", "first", "-n", "second"},
			want: "second",
		},
		"same option twice": {
			args: []string{"--name=first", "--name=second"},
			want: "second",
		},
	}

	for method, parse := range parseMethods {
		for msg, tc := range testCases {
			t.Run(method+"/"+msg, func(t *testing.T) {
				t.Parallel()

				var cfg retryCfg
				og := newRetryGroup(&cfg)

				if err := parse(og, tc.args); err != nil {
					t.Fatalf("og.%s(%v) returns err = %v; want nil", method, tc.args, err)
				}

				if cfg.name != tc.want {
					t.Errorf("og.%s(%v) assigns %q to cfg.name; want %q", method, tc.args, cfg.name, tc.want)
				}
			})
		}
	}
}
