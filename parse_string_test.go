// parse_string_test.go
package opts_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/opts"
)

func TestParseSingleDashStringOption(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		want     string
		args     []string
		postArgs []string
	}{
		"Basic value; single dash": {
			args:     []string{"-file", "test.txt"},
			postArgs: []string{},
			want:     "test.txt",
		},
		"Value with spaces; single dash": {
			args:     []string{"-file", "test file.txt"},
			postArgs: []string{},
			want:     "test file.txt",
		},
		"Empty string; single dash": {
			args:     []string{"-file", ""},
			postArgs: []string{},
			want:     "",
		},
		"Args after value; single dash": {
			args:     []string{"-file", "test.txt", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
			want:     "test.txt",
		},
		"Space separated; double dash": {
			args:     []string{"--file", "test.txt"},
			postArgs: []string{},
			want:     "test.txt",
		},
		"With equals; double dash": {
			args:     []string{"--file=test.txt"},
			postArgs: []string{},
			want:     "test.txt",
		},
		"Value with spaces; double dash": {
			args:     []string{"--file", "test file.txt"},
			postArgs: []string{},
			want:     "test file.txt",
		},
		"Empty string; double dash": {
			args:     []string{"--file", ""},
			postArgs: []string{},
			want:     "",
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

			err := og.Parse(tc.args)
			if err != nil {
				t.Fatalf("after err := og.Parse(%v), err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("after og.Parse(%v), got = %q; want %q", tc.args, got, tc.want)
			}

			postArgs := og.Args()
			if diff := cmp.Diff(tc.postArgs, postArgs); diff != "" {
				t.Errorf("after og.Parse(%v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

func TestParseStringOptionErrors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args []string
	}{
		"No value argument; single dash": {
			args: []string{"-file"},
		},
		"No value argument; double dash": {
			args: []string{"--file"},
		},
		"Equal and empty value; double dash": {
			args: []string{"--file="},
		},
		"Unknown option; single dash": {
			args: []string{"-foo", "bar"},
		},
		"Unknown option; double dash": {
			args: []string{"--foo", "bar"},
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
				t.Errorf("after og.Parse(%v), err == nil; want error", tc.args)
			}
		})
	}
}
