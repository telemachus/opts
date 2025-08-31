package opts_test

import (
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
				t.Fatalf("after og.Parse(%v), err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("after og.Parse(%v), got = %q; want %q", tc.args, got, tc.want)
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
				t.Fatalf("after og.ParseKnown(%v), err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("after og.ParseKnown(%v), got = %q; want %q", tc.args, got, tc.want)
			}

			if diff := cmp.Diff(tc.postArgs, remaining); diff != "" {
				t.Errorf("after og.ParseKnown(%v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

func TestParseStringErrors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		assertErr func(t *testing.T, err error)
		args      []string
	}{
		"No value argument; single dash": {
			args:      []string{"-file"},
			assertErr: checkErrorAs[*opts.MissingValueError],
		},
		"No value argument; double dash": {
			args:      []string{"--file"},
			assertErr: checkErrorAs[*opts.MissingValueError],
		},
		"Equal and empty value; double dash": {
			args:      []string{"--file="},
			assertErr: checkErrorAs[*opts.MissingValueError],
		},
		"Unknown option; single dash": {
			args:      []string{"-foo", "bar"},
			assertErr: checkErrorAs[*opts.UnknownOptionError],
		},
		"Unknown option; double dash": {
			args:      []string{"--foo", "bar"},
			assertErr: checkErrorAs[*opts.UnknownOptionError],
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got string
			og := opts.NewGroup("test-parsing")
			og.String(&got, "file", "whatever")

			err := og.Parse(tc.args)
			if err == nil {
				t.Fatalf("after og.Parse(%v), err == nil; want error", tc.args)
			}

			tc.assertErr(t, err)
		})
	}
}
