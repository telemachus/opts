package opts_test

import (
	"errors"
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
				t.Fatalf("og.Parse(%v) returns err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("og.Parse(%v) assigns %d to got; want %d", tc.args, got, tc.want)
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
				t.Fatalf("og.ParseKnown(%v) returns err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("og.ParseKnown(%v) assigns %d to got; want %d", tc.args, got, tc.want)
			}

			if diff := cmp.Diff(tc.postArgs, remaining); diff != "" {
				t.Errorf("og.ParseKnown(%v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

func TestParseIntSimpleErrors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		errWanted error
		args      []string
	}{
		"Single dash, no value": {
			args:      []string{"-n"},
			errWanted: opts.ErrMissingValue,
		},
		"Double dash, no value": {
			args:      []string{"--number"},
			errWanted: opts.ErrMissingValue,
		},
		"Single dash, equals no value": {
			args:      []string{"-n="},
			errWanted: opts.ErrMissingValue,
		},
		"Double dash, equals no value": {
			args:      []string{"--number="},
			errWanted: opts.ErrMissingValue,
		},
		"Single dash, unknown option": {
			args:      []string{"-q"},
			errWanted: opts.ErrUnknownOption,
		},
		"Double dash, unknown option": {
			args:      []string{"--quiet"},
			errWanted: opts.ErrUnknownOption,
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
			if !errors.Is(err, tc.errWanted) {
				t.Fatalf("og.Parse(%v) returns err == %v; want %v", tc.args, err, tc.errWanted)
			}
		})
	}
}

func TestParseIntErrors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		errWanted error
		args      []string
	}{
		"Single dash, invalid value": {
			args: []string{"-n", "xyz"},
		},
		"Double dash, invalid value": {
			args: []string{"--number", "xyz"},
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
			og := opts.NewGroup("test-parsing")
			og.Int(&got, "n", 0)
			og.Int(&got, "number", 0)

			err := og.Parse(tc.args)
			var ive *opts.InvalidValueError
			if !errors.As(err, &ive) {
				t.Fatalf("Parse(%v) returns %T; want InvalidValueError", tc.args, err)
			}
		})
	}
}
