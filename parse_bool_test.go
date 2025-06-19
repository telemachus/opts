package opts_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/opts"
)

func TestParseBool(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args     []string
		postArgs []string
		want     bool
	}{
		"no args; one dash": {
			args:     []string{"-v"},
			postArgs: []string{},
			want:     true,
		},
		"args after option; single dash": {
			args:     []string{"-v", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
			want:     true,
		},
		"no args; two dashes": {
			args:     []string{"--verbose"},
			postArgs: []string{},
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

			err := og.Parse(tc.args)
			if err != nil {
				t.Fatalf("after err := og.Parse(%v), err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("after og.Parse(%v), got = %t; want %t", tc.args, got, tc.want)
			}

			postArgs := og.Args()
			if diff := cmp.Diff(tc.postArgs, postArgs); diff != "" {
				t.Errorf("after og.Parse(%v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

func TestParseBoolError(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args []string
	}{
		"bool with equals true": {
			args: []string{"--verbose=true"},
		},
		"bool with equals false": {
			args: []string{"--verbose=false"},
		},
		"bool with equals empty": {
			args: []string{"--verbose="},
		},
		"bool with equals random string": {
			args: []string{"--verbose=yadda-yadda"},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got bool
			og := opts.NewGroup("test-parsing")
			og.Bool(&got, "verbose")

			err := og.Parse(tc.args)
			if err == nil {
				t.Fatalf("after err := og.Parse(%v), err == nil; want error", tc.args)
			}
		})
	}
}
