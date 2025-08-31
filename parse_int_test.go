package opts_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/opts"
)

func TestParseInt(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args []string
		want int
	}{
		"Basic value; single dash": {
			args: []string{"-n", "42"},
			want: 42,
		},
		"Negative value; single dash": {
			args: []string{"-n", "-42"},
			want: -42,
		},
		"Space separated; double dash": {
			args: []string{"--number", "42"},
			want: 42,
		},
		"With equals; double dash": {
			args: []string{"--number=42"},
			want: 42,
		},
		"Negative value; double dash": {
			args: []string{"--number", "-42"},
			want: -42,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got int
			og := opts.NewGroup("test-parsing")
			og.Int(&got, "n", 0)
			og.Int(&got, "number", 0)

			err := og.Parse(tc.args)
			if err != nil {
				t.Fatalf("after og.Parse(%v), err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("after og.Parse(%v), got = %d; want %d", tc.args, got, tc.want)
			}
		})
	}
}

func TestParseIntWithRemainingArgs(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args     []string
		postArgs []string
		want     int
	}{
		"Args after value; single dash": {
			args:     []string{"-n", "42", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
			want:     42,
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
			og := opts.NewGroup("test-parsing")
			og.Int(&got, "n", 0)
			og.Int(&got, "number", 0)

			remaining, err := og.ParseKnown(tc.args)
			if err != nil {
				t.Fatalf("after og.ParseKnown(%v), err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("after og.ParseKnown(%v), got = %d; want %d", tc.args, got, tc.want)
			}

			if diff := cmp.Diff(tc.postArgs, remaining); diff != "" {
				t.Errorf("after og.ParseKnown(%v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

func TestParseIntErrors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		assertErr func(t *testing.T, err error)
		args      []string
	}{
		"Single dash, no value": {
			args:      []string{"-n"},
			assertErr: checkErrorAs[*opts.MissingValueError],
		},
		"Double dash, no value": {
			args:      []string{"--number"},
			assertErr: checkErrorAs[*opts.MissingValueError],
		},
		"Double dash, equals no value": {
			args:      []string{"--number="},
			assertErr: checkErrorAs[*opts.MissingValueError],
		},
		"Single dash, invalid value": {
			args:      []string{"-n", "xyz"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Double dash, invalid value": {
			args:      []string{"--number", "xyz"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Double dash, equals invalid": {
			args:      []string{"--number=xyz"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Single dash, float value": {
			args:      []string{"-n", "3.14"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Double dash, float value": {
			args:      []string{"--number=3.14"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Double dash, multiple equals": {
			args:      []string{"--number=42=13"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got int
			og := opts.NewGroup("test-parsing")
			og.Int(&got, "n", 0)
			og.Int(&got, "number", 0)

			err := og.Parse(tc.args)
			if err == nil {
				t.Fatalf("after og.Parse(%v), err == nil; want error", tc.args)
			}

			tc.assertErr(t, err)
		})
	}
}
