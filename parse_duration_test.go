package opts_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/opts"
)

func TestParseDuration(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args []string
		want time.Duration
	}{
		"Basic seconds; single dash": {
			args: []string{"-duration", "10s"},
			want: 10 * time.Second,
		},
		"Zero; single dash": {
			args: []string{"-duration", "0s"},
			want: 0,
		},
		"Minutes; single dash": {
			args: []string{"-duration", "5m"},
			want: 5 * time.Minute,
		},
		"Hours; single dash": {
			args: []string{"-duration", "2h"},
			want: 2 * time.Hour,
		},
		"Complex duration; single dash": {
			args: []string{"-duration", "2h30m"},
			want: 2*time.Hour + 30*time.Minute,
		},
		"Milliseconds; single dash": {
			args: []string{"-duration", "1500ms"},
			want: 1500 * time.Millisecond,
		},
		"Space separated; two dashes": {
			args: []string{"--duration", "10s"},
			want: 10 * time.Second,
		},
		"With equals; two dashes": {
			args: []string{"--duration=10s"},
			want: 10 * time.Second,
		},
		"Zero; two dashes": {
			args: []string{"--duration", "0s"},
			want: 0,
		},
		"Minutes; two dashes": {
			args: []string{"--duration=5m"},
			want: 5 * time.Minute,
		},
		"Hours; two dashes": {
			args: []string{"--duration", "2h"},
			want: 2 * time.Hour,
		},
		"Complex duration; two dashes": {
			args: []string{"--duration=2h30m"},
			want: 2*time.Hour + 30*time.Minute,
		},
		"Milliseconds; two dashes": {
			args: []string{"--duration", "1500ms"},
			want: 1500 * time.Millisecond,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got time.Duration
			og := opts.NewGroup("test-parsing")
			og.Duration(&got, "duration", 0)

			err := og.Parse(tc.args)
			if err != nil {
				t.Fatalf("og.Parse(%v) returns %v as err; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("og.Parse(%v) assigns %v to got; want %v", tc.args, got, tc.want)
			}
		})
	}
}

func TestParseDurationWithRemainingArgs(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args     []string
		postArgs []string
		want     time.Duration
	}{
		"Args after value; single dash": {
			args:     []string{"-duration", "1h", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
			want:     time.Hour,
		},
		"Args after value; two dashes": {
			args:     []string{"--duration", "1h", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
			want:     time.Hour,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got time.Duration
			og := opts.NewGroup("test-parsing")
			og.Duration(&got, "duration", 0)

			remaining, err := og.ParseKnown(tc.args)
			if err != nil {
				t.Fatalf("og.ParseKnown(%v) returns %v as err; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("og.ParseKnown(%v) assigns %v to got; want %v", tc.args, got, tc.want)
			}

			if diff := cmp.Diff(tc.postArgs, remaining); diff != "" {
				t.Errorf("og.ParseKnown(%v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

func TestParseDurationSimpleErrors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		errWanted error
		args      []string
	}{
		"No value; single dash": {
			args:      []string{"-duration"},
			errWanted: opts.ErrMissingValue,
		},
		"No value; double dash": {
			args:      []string{"--duration"},
			errWanted: opts.ErrMissingValue,
		},
		"Equals without value": {
			args:      []string{"--duration="},
			errWanted: opts.ErrMissingValue,
		},
		"Unknown option": {
			args:      []string{"-foo=1h"},
			errWanted: opts.ErrUnknownOption,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got time.Duration
			og := opts.NewGroup("test-parsing")
			og.Duration(&got, "d", 0)
			og.Duration(&got, "duration", 0)

			err := og.Parse(tc.args)
			if !errors.Is(err, tc.errWanted) {
				t.Errorf("og.Parse(%v), got %v; want %v", tc.args, err, tc.errWanted)
			}
		})
	}
}

func TestParseDurationInvalidValueErrors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args []string
	}{
		"Invalid value; single dash": {
			args: []string{"-duration", "xyz"},
		},
		"Invalid value; double dash": {
			args: []string{"--duration", "xyz"},
		},
		"Equals with invalid value": {
			args: []string{"--duration=xyz"},
		},
		"Missing unit; single dash": {
			args: []string{"-duration", "42"},
		},
		"Missing unit; double dash": {
			args: []string{"--duration=42"},
		},
		"Invalid unit": {
			args: []string{"-duration", "42y"},
		},
		"Negative without number": {
			args: []string{"-duration", "-s"},
		},
		"Double dash, multiple equals": {
			args: []string{"--duration=1h=2h"},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got time.Duration
			og := opts.NewGroup("test-parsing")
			og.Duration(&got, "d", 0)
			og.Duration(&got, "duration", 0)

			err := og.Parse(tc.args)
			var ive *opts.InvalidValueError

			if !errors.As(err, &ive) {
				t.Errorf("og.Parse(%v) returns %T; want InvalidValueError", tc.args, err)
			}
		})
	}
}
