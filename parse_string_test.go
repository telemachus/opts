package opts_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/opts"
)

func TestParseString(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		want string
		args []string
	}{
		"Basic value; single dash": {
			args: []string{"-file", "test.txt"},
			want: "test.txt",
		},
		"Value with spaces; single dash": {
			args: []string{"-file", "test file.txt"},
			want: "test file.txt",
		},
		"Empty string; single dash": {
			args: []string{"-file", ""},
			want: "",
		},
		"Space separated; double dash": {
			args: []string{"--file", "test.txt"},
			want: "test.txt",
		},
		"With equals; double dash": {
			args: []string{"--file=test.txt"},
			want: "test.txt",
		},
		"Value with spaces; double dash": {
			args: []string{"--file", "test file.txt"},
			want: "test file.txt",
		},
		"Empty string; double dash": {
			args: []string{"--file", ""},
			want: "",
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got string
			og := opts.NewGroup("test-parsing")
			og.String(&got, "file", "whatever")

			err := og.Parse(tc.args)
			if err != nil {
				t.Fatalf("og.Parse(%v) returns err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("og.Parse(%v) assigns %q to got; want %q", tc.args, got, tc.want)
			}
		})
	}
}

func TestParseStringWithRemainingArgs(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		want     string
		args     []string
		postArgs []string
	}{
		"Args after value; single dash": {
			args:     []string{"-file", "test.txt", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
			want:     "test.txt",
		},
		"Args after value; double dash": {
			args:     []string{"--file", "test.txt", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
			want:     "test.txt",
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got string
			og := opts.NewGroup("test-parsing")
			og.String(&got, "file", "whatever")

			remaining, err := og.ParseKnown(tc.args)
			if err != nil {
				t.Fatalf("og.ParseKnown(%v) returns err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("og.ParseKnown(%v) assigns %q to got; want %q", tc.args, got, tc.want)
			}

			if diff := cmp.Diff(tc.postArgs, remaining); diff != "" {
				t.Errorf("og.ParseKnown(%v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

func TestParseStringSimpleErrors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		errWanted error
		args      []string
	}{
		"No value argument; single dash": {
			args:      []string{"-file"},
			errWanted: opts.ErrMissingValue,
		},
		"No value argument; double dash": {
			args:      []string{"--file"},
			errWanted: opts.ErrMissingValue,
		},
		"Equal and empty value; double dash": {
			args:      []string{"--file="},
			errWanted: opts.ErrMissingValue,
		},
		"Unknown option; single dash": {
			args:      []string{"-foo", "bar"},
			errWanted: opts.ErrUnknownOption,
		},
		"Unknown option; double dash": {
			args:      []string{"--foo", "bar"},
			errWanted: opts.ErrUnknownOption,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got string
			og := opts.NewGroup("test-parsing")
			og.String(&got, "file", "whatever")

			err := og.Parse(tc.args)
			if !errors.Is(err, tc.errWanted) {
				t.Fatalf("og.Parse(%v) returns err == nil; want %v", tc.args, tc.errWanted)
			}
		})
	}
}
