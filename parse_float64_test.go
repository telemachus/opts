package opts_test

import (
	"errors"
	"testing"

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
		"Scientific notation without sign; double dash": {
			args: []string{"--value=1e3"},
			want: 1000.0,
		},
		"No leading digit; single dash": {
			args: []string{"-x", ".5"},
			want: 0.5,
		},
		"Explicit plus sign; double dash": {
			args: []string{"--value", "+3"},
			want: 3.0,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got float64
			og := opts.NewGroup("test-parsing")
			og.Float64(&got, "x", 0.0)
			og.Float64(&got, "value", 0.0)

			err := og.ParseAll(tc.args)
			if err != nil {
				t.Fatalf("og.ParseAll(%v) returns %v as err; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("og.ParseAll(%v) assigns %g to got; want %g", tc.args, got, tc.want)
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

			remaining, err := og.Parse(tc.args)
			if err != nil {
				t.Fatalf("og.Parse(%v) returns %v as err; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("og.Parse(%v) assigns %g to got; want %g", tc.args, got, tc.want)
			}

			checkRemaining(t, "Parse", tc.args, tc.postArgs, remaining)
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

			err := og.ParseAll(tc.args)
			if !errors.Is(err, opts.ErrMissingValue) {
				t.Fatalf("og.ParseAll(%v) returns %v as err; want %v", tc.args, err, opts.ErrMissingValue)
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
		"Single dash, NaN": {
			args: []string{"-x", "NaN"},
		},
		"Double dash, lowercase inf": {
			args: []string{"--value=inf"},
		},
		"Single dash, negative infinity": {
			args: []string{"-x", "-Infinity"},
		},
		"Double dash, hex float": {
			args: []string{"--value", "0x1p4"},
		},
		"Single dash, underscore separator": {
			args: []string{"-x", "1_000.5"},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got float64
			og := opts.NewGroup("test-parsing")
			og.Float64(&got, "x", 0.0)
			og.Float64(&got, "value", 0.0)

			err := og.ParseAll(tc.args)
			if _, ok := errors.AsType[*opts.InvalidValueError](err); !ok {
				t.Fatalf("og.ParseAll(%v) returns %T; want *opts.InvalidValueError", tc.args, err)
			}
		})
	}
}
