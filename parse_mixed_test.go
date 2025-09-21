package opts_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/opts"
)

func TestParseMultipleDifferentOptions(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		wantS string
		args  []string
		wantI int
		wantF float64
		wantB bool
	}{
		"Mixed single dash options": {
			args:  []string{"-v", "-n", "42", "-x", "3.14", "-f", "test.txt"},
			wantB: true,
			wantI: 42,
			wantF: 3.14,
			wantS: "test.txt",
		},
		"Mixed double dash options": {
			args:  []string{"--verbose", "--number", "42", "--value", "3.14", "--file", "test.txt"},
			wantB: true,
			wantI: 42,
			wantF: 3.14,
			wantS: "test.txt",
		},
		"Mixed single and double dash options": {
			args:  []string{"-v", "--number", "42", "-x", "3.14", "--file", "test.txt"},
			wantB: true,
			wantI: 42,
			wantF: 3.14,
			wantS: "test.txt",
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var gotB bool
			var gotI int
			var gotF float64
			var gotS string

			og := opts.NewGroup("test-parsing")
			og.Bool(&gotB, "v")
			og.Bool(&gotB, "verbose")
			og.Int(&gotI, "n", 0)
			og.Int(&gotI, "number", 0)
			og.Float64(&gotF, "x", 0.0)
			og.Float64(&gotF, "value", 0.0)
			og.String(&gotS, "f", "")
			og.String(&gotS, "file", "")

			err := og.Parse(tc.args)
			if err != nil {
				t.Fatalf("og.Parse(%v) returns %v as err; want nil", tc.args, err)
			}

			if gotB != tc.wantB {
				t.Errorf("og.Parse(%v) assigns %t to bool; want %t", tc.args, gotB, tc.wantB)
			}
			if gotI != tc.wantI {
				t.Errorf("og.Parse(%v) assigns %d to int; want %d", tc.args, gotI, tc.wantI)
			}
			if gotF != tc.wantF {
				t.Errorf("og.Parse(%v) assigns %g to float; want %g", tc.args, gotF, tc.wantF)
			}
			if gotS != tc.wantS {
				t.Errorf("og.Parse(%v) assigns %q to string; want %q", tc.args, gotS, tc.wantS)
			}
		})
	}
}

func TestParseMultipleDifferentOptionsWithRemainingArgs(t *testing.T) {
	t.Parallel()

	args := []string{"-v", "-n=42", "--value=3.14", "-f", "test.txt", "foo", "bar"}
	postArgs := []string{"foo", "bar"}

	var gotB bool
	var gotI int
	var gotF float64
	var gotS string

	og := opts.NewGroup("test-parsing")
	og.Bool(&gotB, "v")
	og.Bool(&gotB, "verbose")
	og.Int(&gotI, "n", 0)
	og.Int(&gotI, "number", 0)
	og.Float64(&gotF, "x", 0.0)
	og.Float64(&gotF, "value", 0.0)
	og.String(&gotS, "f", "")
	og.String(&gotS, "file", "")

	remaining, err := og.ParseKnown(args)
	if err != nil {
		t.Fatalf("og.ParseKnown(%v) returns %v as err; want nil", args, err)
	}

	if gotB != true {
		t.Errorf("og.ParseKnown(%v) assigns %t to bool; want %t", args, gotB, true)
	}
	if gotI != 42 {
		t.Errorf("og.ParseKnown(%v) assigns %d to int; want %d", args, gotI, 42)
	}
	if gotF != 3.14 {
		t.Errorf("og.ParseKnown(%v) assigns %g to float; want %g", args, gotF, 3.14)
	}
	if gotS != "test.txt" {
		t.Errorf("og.ParseKnown(%v) assigns %q to string; want %q", args, gotS, "test.txt")
	}

	if diff := cmp.Diff(postArgs, remaining); diff != "" {
		t.Errorf("og.ParseKnown(%v); (-want +got):\n%s", args, diff)
	}
}
