package opts_test

import (
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
				t.Fatalf("after og.Parse(%v), err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("after og.Parse(%v), got = %v; want %v", tc.args, got, tc.want)
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
				t.Fatalf("after og.ParseKnown(%v), err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("after og.ParseKnown(%v), got = %v; want %v", tc.args, got, tc.want)
			}

			if diff := cmp.Diff(tc.postArgs, remaining); diff != "" {
				t.Errorf("after og.ParseKnown(%v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

func TestParseDurationErrors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		assertErr func(t *testing.T, err error)
		args      []string
	}{
		"No value; single dash": {
			args:      []string{"-duration"},
			assertErr: checkErrorAs[*opts.MissingValueError],
		},
		"No value; double dash": {
			args:      []string{"--duration"},
			assertErr: checkErrorAs[*opts.MissingValueError],
		},
		"Invalid value; single dash": {
			args:      []string{"-duration", "xyz"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Invalid value; double dash": {
			args:      []string{"--duration", "xyz"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Equals without value": {
			args:      []string{"--duration="},
			assertErr: checkErrorAs[*opts.MissingValueError],
		},
		"Equals with invalid value": {
			args:      []string{"--duration=xyz"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Missing unit; single dash": {
			args:      []string{"-duration", "42"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Missing unit; double dash": {
			args:      []string{"--duration=42"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Invalid unit": {
			args:      []string{"-duration", "42y"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Negative without number": {
			args:      []string{"-duration", "-s"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Double dash, multiple equals": {
			args:      []string{"--duration=1h=2h"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Unregistered option": {
			args:      []string{"-foo=1h"},
			assertErr: checkErrorAs[*opts.UnknownOptionError],
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
			if err == nil {
				t.Errorf("after og.Parse(%v), err == nil; want error", tc.args)
			}
		})
	}
}
