package opts_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/opts"
)

func TestParseBool(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args []string
		want bool
	}{
		"no args; one dash": {
			args: []string{"-v"},
			want: true,
		},
		"no args; two dashes": {
			args: []string{"--verbose"},
			want: true,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got bool
			og := opts.NewGroup("test-parsing")
			og.Bool(&got, "v")
			og.Bool(&got, "verbose")

			err := og.Parse(tc.args)
			if err != nil {
				t.Fatalf("og.Parse(%v) returns %v as err; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("og.Parse(%v) assigns %t to got; want %t", tc.args, got, tc.want)
			}
		})
	}
}

func TestParseBoolWithRemainingArgs(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args     []string
		postArgs []string
		want     bool
	}{
		"args after option; single dash": {
			args:     []string{"-v", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
			want:     true,
		},
		"args after option; two dashes": {
			args:     []string{"--verbose", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
			want:     true,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got bool
			og := opts.NewGroup("test-parsing")
			og.Bool(&got, "v")
			og.Bool(&got, "verbose")

			remaining, err := og.ParseKnown(tc.args)
			if err != nil {
				t.Fatalf("og.ParseKnown(%v) returns %v as err; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("og.ParseKnown(%v) assigns %t to got; want %t", tc.args, got, tc.want)
			}

			if diff := cmp.Diff(tc.postArgs, remaining); diff != "" {
				t.Errorf("og.ParseKnown(%v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

func TestParseBoolErrors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		errWanted error
		args      []string
	}{
		"bool with equals true": {
			args:      []string{"--verbose=true"},
			errWanted: opts.ErrBooleanWithValue,
		},
		"bool with equals false": {
			args:      []string{"--verbose=false"},
			errWanted: opts.ErrBooleanWithValue,
		},
		"bool with equals empty": {
			args:      []string{"--verbose="},
			errWanted: opts.ErrBooleanWithValue,
		},
		"bool with equals random string": {
			args:      []string{"--verbose=yadda-yadda"},
			errWanted: opts.ErrBooleanWithValue,
		},
		"unknown option": {
			args:      []string{"--foobar"},
			errWanted: opts.ErrUnknownOption,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got bool
			og := opts.NewGroup("test-parsing")
			og.Bool(&got, "verbose")

			err := og.Parse(tc.args)
			if !errors.Is(err, tc.errWanted) {
				t.Errorf("og.Parse(%v) returns %v as err; want %v", tc.args, err, tc.errWanted)
			}

			og2 := opts.NewGroup("test-parsing")
			og2.Bool(&got, "verbose")
			_, err2 := og2.ParseKnown(tc.args)
			if !errors.Is(err2, tc.errWanted) {
				t.Errorf("og.Parse(%v) returns %v as err; want %v", tc.args, err2, tc.errWanted)
			}
		})
	}
}
