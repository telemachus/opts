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
			g := opts.NewGroup("test-parsing")
			g.BoolZero(&got, "v")
			g.BoolZero(&got, "verbose")

			err := g.Parse(tc.args)
			if err != nil {
				t.Fatalf("after err := g.Parse(%+v), err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("after g.Parse(%+v), got = %t; want %t", tc.args, got, tc.want)
			}

			postArgs := g.Args()
			if diff := cmp.Diff(tc.postArgs, postArgs); diff != "" {
				t.Errorf("g.Parse(%+v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}
