// parse_bool_test.go
package opts_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/opts"
)

func TestParseBool(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args     []string
		postArgs []string
		want     bool
	}{
		"no args; one dash": {
			args:     []string{"-v"},
			postArgs: []string{},
			want:     true,
		},
		"args after flag; one dash": {
			args:     []string{"-v", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
			want:     true,
		},
		"no args; two dashes": {
			args:     []string{"--verbose"},
			postArgs: []string{},
			want:     true,
		},
		"args after flag; two dashes": {
			args:     []string{"--verbose", "foo", "bar"},
			postArgs: []string{"foo", "bar"},
			want:     true,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got bool
			g := opts.NewGroup("test-parsing")
			g.Bool(&got, "v")
			g.Bool(&got, "verbose")

			err := g.Parse(tc.args)
			if err != nil {
				t.Fatalf("after err := g.Parse(%+v), err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("after g.Parse(%+v), got = %t; want %t", tc.args, got, tc.want)
			}

			postArgs := g.Args()
			if diff := cmp.Diff(tc.postArgs, postArgs); diff != "" {
				t.Errorf("g.Parse(%+v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

// func TestParseMultipleShortBoolFlags(t *testing.T) {
// 	t.Parallel()
//
// 	testCases := map[string]struct {
// 		args     []string
// 		postArgs []string
// 		wantA    bool
// 		wantB    bool
// 		wantC    bool
// 	}{
// 		"Three flags": {
// 			args:     []string{"-abc", "foo", "bar"},
// 			postArgs: []string{"foo", "bar"},
// 			wantA:    true,
// 			wantB:    true,
// 			wantC:    true,
// 		},
// 		"Two flags": {
// 			args:     []string{"-ab", "foo", "bar"},
// 			postArgs: []string{"foo", "bar"},
// 			wantA:    true,
// 			wantB:    true,
// 			wantC:    false,
// 		},
// 	}
//
// 	for msg, tc := range testCases {
// 		t.Run(msg, func(t *testing.T) {
// 			t.Parallel()
//
// 			var gotA, gotB, gotC bool
// 			fs := flagx.NewFlagSet("test-parsing")
// 			fs.BoolVar(&gotA, 'a', "", false)
// 			fs.BoolVar(&gotB, 'b', "", false)
// 			fs.BoolVar(&gotC, 'c', "", false)
//
// 			err := fs.Parse(tc.args)
// 			if err != nil {
// 				t.Fatalf("after err := fs.Parse(%+v), err == %v; want nil", tc.args, err)
// 			}
//
// 			if gotA != tc.wantA {
// 				t.Errorf("after fs.Parse(%+v), -a = %t; want %t", tc.args, gotA, tc.wantA)
// 			}
// 			if gotB != tc.wantB {
// 				t.Errorf("after fs.Parse(%+v), -b = %t; want %t", tc.args, gotB, tc.wantB)
// 			}
// 			if gotC != tc.wantC {
// 				t.Errorf("after fs.Parse(%+v), -c = %t; want %t", tc.args, gotC, tc.wantC)
// 			}
//
// 			postArgs := fs.Args()
// 			if diff := cmp.Diff(tc.postArgs, postArgs); diff != "" {
// 				t.Errorf("fs.Parse(%+v); (-want +got):\n%s", tc.args, diff)
// 			}
// 		})
// 	}
// }
//
// func TestParseBoolFlagErrors(t *testing.T) {
// 	t.Parallel()
//
// 	testCases := map[string]struct {
// 		args []string
// 	}{
// 		"Short flag with equal and valid value": {
// 			args: []string{"-v=true"},
// 		},
// 		"Short flag with equal and empty value": {
// 			args: []string{"-v="},
// 		},
// 		"Short flag with equal and invalid value": {
// 			args: []string{"-v=bar"},
// 		},
// 		"Long flag with equal and empty value": {
// 			args: []string{"--verbose="},
// 		},
// 		"Long flag with equal and invalid value": {
// 			args: []string{"--verbose=notabool"},
// 		},
// 	}
//
// 	for msg, tc := range testCases {
// 		t.Run(msg, func(t *testing.T) {
// 			t.Parallel()
//
// 			var got bool
// 			fs := flagx.NewFlagSet("test-parsing")
// 			fs.BoolVar(&got, 'v', "verbose", false)
//
// 			err := fs.Parse(tc.args)
// 			if err == nil {
// 				t.Errorf("after fs.Parse(%+v), err == nil; want error", tc.args)
// 			}
// 		})
// 	}
// }
