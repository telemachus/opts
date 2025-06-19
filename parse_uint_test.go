package opts_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/opts"
)

func TestParseSingleDashUintOption(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args     []string
		postArgs []string
		want     uint
	}{
		"Basic value; single dash": {
			args:     []string{"-n", "42"},
			postArgs: []string{},
			want:     42,
		},
		"Zero; single dash": {
			args:     []string{"-n", "0"},
			postArgs: []string{},
			want:     0,
		},
		"Hex value; single dash": {
			args:     []string{"-n", "0xff"},
			postArgs: []string{},
			want:     255,
		},
		"Octal value; single dash": {
			args:     []string{"-n", "0644"},
			postArgs: []string{},
			want:     420,
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
		"Zero; double dash": {
			args:     []string{"--number", "0"},
			postArgs: []string{},
			want:     0,
		},
		"Hex value; double dash": {
			args:     []string{"--number=0xff"},
			postArgs: []string{},
			want:     255,
		},
		"Octal value; double dash": {
			args:     []string{"--number=0644"},
			postArgs: []string{},
			want:     420,
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

			var got uint
			og := opts.NewGroup("test-parsing")
			og.Uint(&got, "n", 0)
			og.Uint(&got, "number", 0)

			err := og.Parse(tc.args)
			if err != nil {
				t.Fatalf("after err := og.Parse(%v), err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("after og.Parse(%v), got = %d; want %d", tc.args, got, tc.want)
			}

			postArgs := og.Args()
			if diff := cmp.Diff(tc.postArgs, postArgs); diff != "" {
				t.Errorf("after og.Parse(%v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

func TestParseUintErrors(t *testing.T) {
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
		"Single dash, negative value": {
			args: []string{"-n", "-42"},
		},
		"Double dash, negative value": {
			args: []string{"--number=-42"},
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

			var got uint
			og := opts.NewGroup("test-parsing")
			og.Uint(&got, "n", 0)
			og.Uint(&got, "number", 0)

			err := og.Parse(tc.args)
			if err == nil {
				t.Errorf("after og.Parse(%v), err == nil; want error", tc.args)
			}
		})
	}
}
