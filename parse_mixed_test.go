package opts_test

import (
	"testing"

	"github.com/telemachus/opts"
)

type mixedCfg struct {
	s string
	i int
	f float64
	b bool
}

func TestParseMultipleDifferentOptions(t *testing.T) {
	t.Parallel()

	want := mixedCfg{s: "test.txt", i: 42, f: 3.14, b: true}

	testCases := map[string]struct {
		args []string
	}{
		"Mixed single dash options": {
			args: []string{"-v", "-n", "42", "-x", "3.14", "-f", "test.txt"},
		},
		"Mixed double dash options": {
			args: []string{"--verbose", "--number", "42", "--value", "3.14", "--file", "test.txt"},
		},
		"Mixed single and double dash options": {
			args: []string{"-v", "--number", "42", "-x", "3.14", "--file", "test.txt"},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got mixedCfg

			og := opts.NewGroup("test-parsing")
			og.Bool(&got.b, "v")
			og.Bool(&got.b, "verbose")
			og.Int(&got.i, "n", 0)
			og.Int(&got.i, "number", 0)
			og.Float64(&got.f, "x", 0.0)
			og.Float64(&got.f, "value", 0.0)
			og.String(&got.s, "f", "")
			og.String(&got.s, "file", "")

			err := og.ParseAll(tc.args)
			if err != nil {
				t.Fatalf("og.ParseAll(%v) returns %v as err; want nil", tc.args, err)
			}

			if got != want {
				t.Errorf("og.ParseAll(%v) assigns %+v; want %+v", tc.args, got, want)
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

	remaining, err := og.Parse(args)
	if err != nil {
		t.Fatalf("og.Parse(%v) returns %v as err; want nil", args, err)
	}

	if !gotB {
		t.Errorf("og.Parse(%v) assigns %t to bool; want %t", args, gotB, true)
	}
	if gotI != 42 {
		t.Errorf("og.Parse(%v) assigns %d to int; want %d", args, gotI, 42)
	}
	if gotF != 3.14 {
		t.Errorf("og.Parse(%v) assigns %g to float; want %g", args, gotF, 3.14)
	}
	if gotS != "test.txt" {
		t.Errorf("og.Parse(%v) assigns %q to string; want %q", args, gotS, "test.txt")
	}

	checkRemaining(t, "Parse", args, postArgs, remaining)
}
