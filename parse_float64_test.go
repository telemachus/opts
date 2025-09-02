package opts_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/opts"
)

func TestParseFloat64(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args []string
		want float64
	}{
		"Integer value; single dash": {
			args: []string{"-x", "42"},
			want: 42.0,
		},
		"Decimal value; single dash": {
			args: []string{"-x", "3.14"},
			want: 3.14,
		},
		"Scientific notation; single dash": {
			args: []string{"-x", "1e-2"},
			want: 0.01,
		},
		"Negative value; single dash": {
			args: []string{"-x", "-3.14"},
			want: -3.14,
		},
		"Space separated decimal; double dash": {
			args: []string{"--value", "3.14"},
			want: 3.14,
		},
		"With equals decimal; double dash": {
			args: []string{"--value=3.14"},
			want: 3.14,
		},
		"Scientific notation; double dash": {
			args: []string{"--value=1e-2"},
			want: 0.01,
		},
		"Integer value; double dash": {
			args: []string{"--value", "42"},
			want: 42.0,
		},
		"Negative value; double dash": {
			args: []string{"--value", "-3.14"},
			want: -3.14,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got float64
			og := opts.NewGroup("test-parsing")
			og.Float64(&got, "x", 0.0)
			og.Float64(&got, "value", 0.0)

			err := og.Parse(tc.args)
			if err != nil {
				t.Fatalf("og.Parse(%v) returns %v as err; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("og.Parse(%v) assigns %g to got; want %g", tc.args, got, tc.want)
			}
		})
	}
}

func TestParseFloat64WithRemainingArgs(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args     []string
		postArgs []string
		want     float64
	}{
		"Args after value; single dash": {
			args:     []string{"-x", "3.14", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
			want:     3.14,
		},
		"Args after value; double dash": {
			args:     []string{"--value", "3.14", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
			want:     3.14,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got float64
			og := opts.NewGroup("test-parsing")
			og.Float64(&got, "x", 0.0)
			og.Float64(&got, "value", 0.0)

			remaining, err := og.ParseKnown(tc.args)
			if err != nil {
				t.Fatalf("og.ParseKnown(%v) returns %v as err; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("og.ParseKnown(%v) assigns %g to got; want %g", tc.args, got, tc.want)
			}

			if diff := cmp.Diff(tc.postArgs, remaining); diff != "" {
				t.Errorf("og.ParseKnown(%v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

func TestParseFloat64ErrMissingValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		errWanted error
		args      []string
	}{
		"Single dash, no value": {
			args:      []string{"-x"},
			errWanted: opts.ErrMissingValue,
		},
		"Double dash, no value": {
			args:      []string{"--value"},
			errWanted: opts.ErrMissingValue,
		},
		"Double dash, equals no value": {
			args:      []string{"--value="},
			errWanted: opts.ErrMissingValue,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got float64
			og := opts.NewGroup("test-parsing")
			og.Float64(&got, "x", 0.0)
			og.Float64(&got, "value", 0.0)

			err := og.Parse(tc.args)
			if !errors.Is(err, opts.ErrMissingValue) {
				t.Fatalf("og.Parse(%v) returns %v as err; want %v", tc.args, err, opts.ErrMissingValue)
			}
		})
	}
}

func TestParseFloat64InvalidValueError(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args []string
	}{
		"Single dash, invalid value": {
			args: []string{"-x", "xyz"},
		},
		"Double dash, invalid value": {
			args: []string{"--value", "xyz"},
		},
		"Double dash, equals invalid": {
			args: []string{"--value=xyz"},
		},
		"Double dash, invalid scientific": {
			args: []string{"--value=1e"},
		},
		"Double dash, multiple equals": {
			args: []string{"--value=3.14=2.718"},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got float64
			og := opts.NewGroup("test-parsing")
			og.Float64(&got, "x", 0.0)
			og.Float64(&got, "value", 0.0)

			err := og.Parse(tc.args)
			var ive *opts.InvalidValueError

			if !errors.As(err, &ive) {
				t.Fatalf("og.Parse(%v) returns %T; want *opts.InvalidValueError", tc.args, err)
			}
		})
	}
}
