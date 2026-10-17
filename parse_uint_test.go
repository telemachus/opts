package opts_test

import (
	"errors"
	"testing"

	"github.com/telemachus/opts"
)

func TestParseUint(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args []string
		want uint
	}{
		"Basic value; single dash": {
			args: []string{"-n", "42"},
			want: 42,
		},
		"Zero; single dash": {
			args: []string{"-n", "0"},
			want: 0,
		},
		"Leading zero is decimal; single dash": {
			args: []string{"-n", "0644"},
			want: 644,
		},
		"Space separated; double dash": {
			args: []string{"--number", "42"},
			want: 42,
		},
		"With equals; double dash": {
			args: []string{"--number=42"},
			want: 42,
		},
		"Zero; double dash": {
			args: []string{"--number", "0"},
			want: 0,
		},
		"Leading zero is decimal; double dash": {
			args: []string{"--number=010"},
			want: 10,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got uint
			og := opts.NewGroup("test-parsing")
			og.Uint(&got, "n", 0)
			og.Uint(&got, "number", 0)

			err := og.ParseAll(tc.args)
			if err != nil {
				t.Fatalf("og.ParseAll(%v) returns %v as err; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("og.ParseAll(%v) assigns %d to got; want %d", tc.args, got, tc.want)
			}
		})
	}
}

func TestParseUintWithRemainingArgs(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args     []string
		postArgs []string
		want     uint
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

			var got uint
			og := opts.NewGroup("test-parsing")
			og.Uint(&got, "n", 0)
			og.Uint(&got, "number", 0)

			remaining, err := og.Parse(tc.args)
			if err != nil {
				t.Fatalf("og.Parse(%v) returns %v as err; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("og.Parse(%v) assigns %d to got; want %d", tc.args, got, tc.want)
			}

			checkRemaining(t, "Parse", tc.args, tc.postArgs, remaining)
		})
	}
}

func TestParseUintSimpleErrors(t *testing.T) {
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
			args:      []string{"-number="},
			errWanted: opts.ErrMissingValue,
		},
		"Double dash, equals no value": {
			args:      []string{"--number="},
			errWanted: opts.ErrMissingValue,
		},
		"Single dash, equals unknown option": {
			args:      []string{"-foobar="},
			errWanted: opts.ErrUnknownOption,
		},
		"Double dash, unknown option": {
			args:      []string{"--foobar"},
			errWanted: opts.ErrUnknownOption,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got uint
			og := opts.NewGroup("test-parsing")
			og.Uint(&got, "n", 0)
			og.Uint(&got, "number", 0)

			err := og.ParseAll(tc.args)
			if !errors.Is(err, tc.errWanted) {
				t.Fatalf("og.ParseAll(%v) returns %v as err; want %v", tc.args, err, tc.errWanted)
			}
		})
	}
}

func TestParseUintInvalidValueError(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args []string
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
		"Single dash, hex value": {
			args: []string{"-n", "0xff"},
		},
		"Double dash, octal prefix": {
			args: []string{"--number=0o17"},
		},
		"Double dash, binary prefix": {
			args: []string{"--number=0b1"},
		},
		"Single dash, underscore separator": {
			args: []string{"-n", "1_000"},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got uint
			og := opts.NewGroup("test-parsing")
			og.Uint(&got, "n", 0)
			og.Uint(&got, "number", 0)

			err := og.ParseAll(tc.args)
			if _, ok := errors.AsType[*opts.InvalidValueError](err); !ok {
				t.Fatalf("ParseAll(%v) returns %T; want InvalidValueError", tc.args, err)
			}
		})
	}
}
