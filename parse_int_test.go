package opts_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/opts"
)

func TestParseInt(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args     []string
		postArgs []string
		want     int
	}{
		"Basic value; single dash": {
			args:     []string{"-n", "42"},
			postArgs: []string{},
			want:     42,
		},
		"Negative value; single dash": {
			args:     []string{"-n", "-42"},
			postArgs: []string{},
			want:     -42,
		},
		"Args after value; single dash": {
			args:     []string{"-n", "42", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
			want:     42,
		},
		"Space separated; double dash": {
			args:     []string{"--number", "42"},
			postArgs: []string{},
			want:     42,
		},
		"With equals; double dash": {
			args:     []string{"--number=42"},
			postArgs: []string{},
			want:     42,
		},
		"Negative value; double dash": {
			args:     []string{"--number", "-42"},
			postArgs: []string{},
			want:     -42,
		},
		"Args after value; double dash": {
			args:     []string{"--number", "42", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
			want:     42,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got int
			g := opts.NewGroup("test-parsing")
			g.Int(&got, "n", 0)
			g.Int(&got, "number", 0)

			err := g.Parse(tc.args)
			if err != nil {
				t.Fatalf("after err := g.Parse(%+v), err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("after g.Parse(%+v), got = %d; want %d", tc.args, got, tc.want)
			}

			postArgs := g.Args()
			if diff := cmp.Diff(tc.postArgs, postArgs); diff != "" {
				t.Errorf("g.Parse(%+v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

func TestParseIntErrors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args []string
	}{
		"Single dash, no value": {
			args: []string{"-n"},
		},
		"Double dash, no value": {
			args: []string{"--number"},
		},
		"Single dash, invalid value": {
			args: []string{"-n", "xyz"},
		},
		"Double dash, invalid value": {
			args: []string{"--number", "xyz"},
		},
		"Double dash, equals no value": {
			args: []string{"--number="},
		},
		"Double dash, equals invalid": {
			args: []string{"--number=xyz"},
		},
		"Single dash, float value": {
			args: []string{"-n", "3.14"},
		},
		"Double dash, float value": {
			args: []string{"--number=3.14"},
		},
		"Double dash, multiple equals": {
			args: []string{"--number=42=13"},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got int
			g := opts.NewGroup("test-parsing")
			g.Int(&got, "n", 0)
			g.Int(&got, "number", 0)

			err := g.Parse(tc.args)
			if err == nil {
				t.Errorf("after g.Parse(%+v), err == nil; want error", tc.args)
			}
		})
	}
}
