package opts_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/opts"
)

func TestParseNoOptions(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args     []string
		postArgs []string
	}{
		"Empty args": {args: []string{}, postArgs: []string{}},
		"-- should not be in fs.Args()": {
			args:     []string{"--", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
		},
		"- should be in fs.Args()": {
			args:     []string{"-", "foo", "bar"},
			postArgs: []string{"-", "foo", "bar"},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			g := opts.NewGroup("test-parsing")

			err := g.Parse(tc.args)
			if err != nil {
				t.Fatalf("after err := g.Parse(%+v), err == %v; want nil", tc.args, err)
			}

			postArgs := g.Args()
			if diff := cmp.Diff(tc.postArgs, postArgs); diff != "" {
				t.Errorf("g.Parse(%+v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

func TestParseOptions(t *testing.T) {
	t.Parallel()

	args := []string{"-V", "--name", "foobar"}
	cfg := struct {
		name    string
		verbose bool
	}{}

	g := opts.NewGroup("test-parsing")
	g.BoolZero(&cfg.verbose, "V")
	g.String(&cfg.name, "name", "default")

	err := g.Parse(args)
	if err != nil {
		t.Errorf("after g.Parse(%+v), err != %v; want nil", args, err)
	}

	if cfg.verbose != true {
		t.Errorf("after g.Parse(%+v), cfg.verbose = %v; want true", args, cfg.verbose)
	}
	if cfg.name != args[2] {
		t.Errorf("after g.Parse(%+v), cfg.name = %v; want %q", args, cfg.name, args[2])
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
			g := opts.NewGroup("test-parsing")
			err := g.Parse(tc.args)
			if err == nil {
				t.Errorf("after g.Parse(%+v), err == nil; want error", tc.args)
			}
		})
	}
}
