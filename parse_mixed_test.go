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
				t.Fatalf("after og.Parse(%v), err == %v; want nil", tc.args, err)
			}

			if gotB != tc.wantB {
				t.Errorf("after og.Parse(%v), bool = %t; want %t", tc.args, gotB, tc.wantB)
			}
			if gotI != tc.wantI {
				t.Errorf("after og.Parse(%v), int = %d; want %d", tc.args, gotI, tc.wantI)
			}
			if gotF != tc.wantF {
				t.Errorf("after og.Parse(%v), float = %g; want %g", tc.args, gotF, tc.wantF)
			}
			if gotS != tc.wantS {
				t.Errorf("after og.Parse(%v), string = %q; want %q", tc.args, gotS, tc.wantS)
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
		t.Fatalf("after og.ParseKnown(%v), err == %v; want nil", args, err)
	}

	if gotB != true {
		t.Errorf("after og.ParseKnown(%v), bool = %t; want %t", args, gotB, true)
	}
	if gotI != 42 {
		t.Errorf("after og.ParseKnown(%v), int = %d; want %d", args, gotI, 42)
	}
	if gotF != 3.14 {
		t.Errorf("after og.ParseKnown(%v), float = %g; want %g", args, gotF, 3.14)
	}
	if gotS != "test.txt" {
		t.Errorf("after og.ParseKnown(%v), string = %q; want %q", args, gotS, "test.txt")
	}

	if diff := cmp.Diff(postArgs, remaining); diff != "" {
		t.Errorf("after og.ParseKnown(%v); (-want +got):\n%s", args, diff)
	}
}
