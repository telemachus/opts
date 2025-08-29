package opts_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/opts"
)

func TestParseBool(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		parseFunc func(*opts.Group, []string) ([]string, error)
		args      []string
		postArgs  []string
		want      bool
	}{
		"no args; one dash": {
			args:      []string{"-v"},
			postArgs:  []string{},
			want:      true,
			parseFunc: (*opts.Group).Parse,
		},
		"args after option; single dash": {
			args:      []string{"-v", "foo", "bar"},
			postArgs:  []string{"foo", "bar"},
			want:      true,
			parseFunc: (*opts.Group).ParseKnown,
		},
		"no args; two dashes": {
			args:      []string{"--verbose"},
			postArgs:  []string{},
			want:      true,
			parseFunc: (*opts.Group).Parse,
		},
		"args after option; two dashes": {
			args:      []string{"--verbose", "foo", "bar"},
			postArgs:  []string{"foo", "bar"},
			want:      true,
			parseFunc: (*opts.Group).ParseKnown,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got bool
			og := opts.NewGroup("test-parsing")
			og.Bool(&got, "v")
			og.Bool(&got, "verbose")

			remaining, err := tc.parseFunc(og, tc.args)
			if err != nil {
				t.Fatalf("after remaining, err := parseFunc(%v), err == %v; want nil", tc.args, err)
			}

			if got != tc.want {
				t.Errorf("after parseFunc(%v), got = %t; want %t", tc.args, got, tc.want)
			}

			if diff := cmp.Diff(tc.postArgs, remaining); diff != "" {
				t.Errorf("after parseFunc(%v); (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

func TestParseBoolError(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		parseFunc func(*opts.Group, []string) ([]string, error)
		args      []string
	}{
		"bool with equals true": {
			args:      []string{"--verbose=true"},
			parseFunc: (*opts.Group).Parse,
		},
		"bool with equals false": {
			args:      []string{"--verbose=false"},
			parseFunc: (*opts.Group).ParseKnown,
		},
		"bool with equals empty": {
			args:      []string{"--verbose="},
			parseFunc: (*opts.Group).Parse,
		},
		"bool with equals random string": {
			args:      []string{"--verbose=yadda-yadda"},
			parseFunc: (*opts.Group).ParseKnown,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var got bool
			og := opts.NewGroup("test-parsing")
			og.Bool(&got, "verbose")

			_, err := tc.parseFunc(og, tc.args)
			if err == nil {
				t.Fatalf("after parseFunc(%v), err == nil; want error", tc.args)
			}
		})
	}
}
