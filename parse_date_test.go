package opts_test

import (
	"testing"

	"cloud.google.com/go/civil"
	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/opts"
)

func TestParseDate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args []string
		want civil.Date
	}{
		"Basic date; single dash": {
			args: []string{"-d", "2025-01-15"},
			want: civil.Date{Year: 2025, Month: 1, Day: 15},
		},
		"Basic date; double dash": {
			args: []string{"--date", "2025-01-15"},
			want: civil.Date{Year: 2025, Month: 1, Day: 15},
		},
		"With equals; double dash": {
			args: []string{"--date=2025-01-15"},
			want: civil.Date{Year: 2025, Month: 1, Day: 15},
		},
		"Leap year date; single dash": {
			args: []string{"-d", "2024-02-29"},
			want: civil.Date{Year: 2024, Month: 2, Day: 29},
		},
		"End of year; double dash": {
			args: []string{"--date", "2025-12-31"},
			want: civil.Date{Year: 2025, Month: 12, Day: 31},
		},
		"Beginning of year; double dash": {
			args: []string{"--date=2025-01-01"},
			want: civil.Date{Year: 2025, Month: 1, Day: 1},
		},
		"Historical date": {
			args: []string{"--date", "1596-03-31"},
			want: civil.Date{Year: 1596, Month: 3, Day: 31},
		},
		"Future date": {
			args: []string{"--date=2099-12-25"},
			want: civil.Date{Year: 2099, Month: 12, Day: 25},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got civil.Date
			og := opts.NewGroup("test-parsing")
			og.Date(&got, "d", civil.Date{})
			og.Date(&got, "date", civil.Date{})

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

func TestParseDateWithRemainingArgs(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args     []string
		postArgs []string
		want     civil.Date
	}{
		"Args after value; single dash": {
			args:     []string{"-d", "2025-06-15", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
			want:     civil.Date{Year: 2025, Month: 6, Day: 15},
		},
		"Args after value; double dash": {
			args:     []string{"--date", "2025-06-15", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
			want:     civil.Date{Year: 2025, Month: 6, Day: 15},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got civil.Date
			og := opts.NewGroup("test-parsing")
			og.Date(&got, "d", civil.Date{})
			og.Date(&got, "date", civil.Date{})

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

func TestParseDateErrors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		assertErr func(t *testing.T, err error)
		args      []string
	}{
		"Single dash, no value": {
			args:      []string{"-d"},
			assertErr: checkErrorAs[*opts.MissingValueError],
		},
		"Double dash, no value": {
			args:      []string{"--date"},
			assertErr: checkErrorAs[*opts.MissingValueError],
		},
		"Double dash, equals no value": {
			args:      []string{"--date="},
			assertErr: checkErrorAs[*opts.MissingValueError],
		},
		"Invalid format; US style": {
			args:      []string{"--date", "01/15/2025"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Invalid format; European style": {
			args:      []string{"--date", "15/01/2025"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Invalid format; dots": {
			args:      []string{"--date", "2025.01.15"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Invalid format; no separators": {
			args:      []string{"--date", "20250115"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Invalid format; reversed": {
			args:      []string{"--date", "15-01-2025"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Invalid date; February 30th": {
			args:      []string{"--date", "2025-02-30"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Invalid date; month 13": {
			args:      []string{"--date", "2025-13-01"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Invalid date; day 32": {
			args:      []string{"--date", "2025-01-32"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Invalid date; (false) leap year": {
			args:      []string{"--date", "2023-02-29"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Invalid date; April 31st": {
			args:      []string{"--date", "2025-04-31"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Invalid format; time included": {
			args:      []string{"--date", "2025-01-15T10:30:00"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Invalid format; partial date": {
			args:      []string{"--date", "2025-01"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Invalid format; text": {
			args:      []string{"--date", "January 15, 2025"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Invalid format; random text": {
			args:      []string{"--date", "xyz"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Double dash, multiple equals": {
			args:      []string{"--date=2025-01-15=2025-01-16"},
			assertErr: checkErrorAs[*opts.InvalidValueError],
		},
		"Unknown option": {
			args:      []string{"--unknown", "2025-01-15"},
			assertErr: checkErrorAs[*opts.UnknownOptionError],
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got civil.Date
			og := opts.NewGroup("test-parsing")
			og.Date(&got, "d", civil.Date{})
			og.Date(&got, "date", civil.Date{})

			err := og.Parse(tc.args)
			if err == nil {
				t.Fatalf("after og.Parse(%v), err == nil; want error", tc.args)
			}

			tc.assertErr(t, err)
		})
	}
}
