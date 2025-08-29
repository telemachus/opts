package opts_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/opts"
)

func TestParseFloat64(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		parseFunc func(*opts.Group, []string) ([]string, error)
		args      []string
		postArgs  []string
		want      float64
	}{
		"Integer value; single dash": {
			args:      []string{"-x", "42"},
			postArgs:  []string{},
			want:      42.0,
			parseFunc: (*opts.Group).Parse,
		},
		"Decimal value; single dash": {
			args:      []string{"-x", "3.14"},
			postArgs:  []string{},
			want:      3.14,
			parseFunc: (*opts.Group).Parse,
		},
		"Scientific notation; single dash": {
			args:      []string{"-x", "1e-2"},
			postArgs:  []string{},
			want:      0.01,
			parseFunc: (*opts.Group).Parse,
		},
		"Negative value; single dash": {
			args:      []string{"-x", "-3.14"},
			postArgs:  []string{},
			want:      -3.14,
			parseFunc: (*opts.Group).Parse,
		},
		"Args after value; single dash": {
			args:      []string{"-x", "3.14", "foo", "bar"},
			postArgs:  []string{"foo", "bar"},
			want:      3.14,
			parseFunc: (*opts.Group).ParseKnown,
		},
		"Space separated decimal; double dash": {
			args:      []string{"--value", "3.14"},
			postArgs:  []string{},
			want:      3.14,
			parseFunc: (*opts.Group).Parse,
		},
		"With equals decimal; double dash": {
			args:      []string{"--value=3.14"},
			postArgs:  []string{},
			want:      3.14,
			parseFunc: (*opts.Group).Parse,
		},
		"Scientific notation; double dash": {
			args:      []string{"--value=1e-2"},
			postArgs:  []string{},
			want:      0.01,
			parseFunc: (*opts.Group).Parse,
		},
		"Integer value; double dash": {
			args:      []string{"--value", "42"},
			postArgs:  []string{},
			want:      42.0,
			parseFunc: (*opts.Group).Parse,
		},
		"Negative value; double dash": {
			args:      []string{"--value", "-3.14"},
			postArgs:  []string{},
			want:      -3.14,
			parseFunc: (*opts.Group).Parse,
		},
		"Args after value; double dash": {
			args:      []string{"--value", "3.14", "foo", "bar"},
			postArgs:  []string{"foo", "bar"},
			want:      3.14,
			parseFunc: (*opts.Group).ParseKnown,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got float64
			og := opts.NewGroup("test-parsing")
			og.Float64(&got, "x", 0.0)
			og.Float64(&got, "value", 0.0)

			postArgs, err := tc.parseFunc(og, tc.args)
			if err != nil {
				t.Fatalf("after err := og.Parse(%v), err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("after og.Parse(%v), got = %g; want %g", tc.args, got, tc.want)
			}

			if diff := cmp.Diff(tc.postArgs, postArgs); diff != "" {
				t.Errorf("after og.Parse(%v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

func TestParseFloat64Errors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args []string
	}{
		"Single dash, no value": {
			args: []string{"-x"},
		},
		"Double dash, no value": {
			args: []string{"--value"},
		},
		"Single dash, invalid value": {
			args: []string{"-x", "xyz"},
		},
		"Double dash, invalid value": {
			args: []string{"--value", "xyz"},
		},
		"Double dash, equals no value": {
			args: []string{"--value="},
		},
		"Double dash, equals invalid": {
			args: []string{"--value=xyz"},
		},
		"Double dash, invalid scientific notation": {
			args: []string{"--value=1e"},
		},
		"Double dash, multiple equals": {
			args: []string{"--value=3.14=2.718"},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got float64
			og := opts.NewGroup("test-parsing")
			og.Float64(&got, "x", 0.0)
			og.Float64(&got, "value", 0.0)

			_, err := og.Parse(tc.args)
			if err == nil {
				t.Errorf("after og.Parse(%v), err == nil; want error", tc.args)
			}
		})
	}
}
