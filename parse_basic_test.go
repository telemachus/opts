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

			og := opts.NewGroup("test-parsing")

			err := og.Parse(tc.args)
			if err != nil {
				t.Fatalf("after err := og.Parse(%v), err == %v; want nil", tc.args, err)
			}

			postArgs := og.Args()
			if diff := cmp.Diff(tc.postArgs, postArgs); diff != "" {
				t.Errorf("after og.Parse(%v); (-want +got):\n%s", tc.args, diff)
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

	og := opts.NewGroup("test-parsing")
	og.Bool(&cfg.verbose, "V")
	og.String(&cfg.name, "name", "default")

	err := og.Parse(args)
	if err != nil {
		t.Errorf("after og.Parse(%v), err != %v; want nil", args, err)
	}

	if cfg.verbose != true {
		t.Errorf("after og.Parse(%v), cfg.verbose = %v; want true", args, cfg.verbose)
	}
	if cfg.name != args[2] {
		t.Errorf("after og.Parse(%v), cfg.name = %v; want %q", args, cfg.name, args[2])
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
				t.Errorf("after og.Parse(%v), err == nil; want error", tc.args)
			}
		})
	}
}

func TestReparseAfterFailure(t *testing.T) {
	t.Parallel()

	og := opts.NewGroup("test")
	var s string
	defValue := "default"
	og.String(&s, "name", defValue)

	badArgs := []string{"--unknown", "value"}
	goodArgs := []string{"--name", "success", "arg1", "arg2"}

	err := og.Parse(badArgs)
	if err == nil {
		t.Fatalf("after og.Parse(%v), err == nil; want error", badArgs)
	}

	if s != defValue {
		t.Errorf("after failed parse, s == %q; want %q", s, defValue)
	}
	if len(og.Args()) != 0 {
		t.Errorf("after failed parse, og.Args() == %v; want empty slice", og.Args())
	}

	err = og.Parse(goodArgs)
	if err != nil {
		t.Fatalf("after og.parse(%v), err == %v; want nil", goodArgs, err)
	}

	if s != goodArgs[1] {
		t.Errorf("after og.parse(%v), s == %q; want %q", goodArgs, s, goodArgs[1])
	}

	postArgs := og.Args()
	if diff := cmp.Diff(goodArgs[2:], postArgs); diff != "" {
		t.Errorf("after og.Parse(%v); (-want +got):\n%s", goodArgs, diff)
	}
}
