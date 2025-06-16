package opts

import (
	"testing"
	"time"
)

func TestNewGroup(t *testing.T) {
	t.Parallel()

	name := "test-opts"
	g := NewGroup(name)
	if g == nil {
		t.Fatalf("NewGroup(%q) returned nil", name)
	}

	if g.Name() != name {
		t.Errorf("g.Name(%q) == %q; want %q", name, g.Name(), name)
	}
}

func TestOptRegistrationValid(t *testing.T) {
	t.Parallel()

	name := "verbose"
	og := NewGroup("test-optiongroup")
	var got bool
	og.Bool(&got, name)

	if opt := og.opts[name]; opt == nil {
		t.Errorf("option --%s not registered", name)
	}
}

func TestDuplicateOptRegistration(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		first  func(*Group)
		second func(*Group)
	}{
		"duplicate bool": {
			first: func(og *Group) {
				var b bool
				og.Bool(&b, "verbose")
			},
			second: func(og *Group) {
				var b bool
				og.Bool(&b, "verbose")
			},
		},
		"duplicate duration": {
			first: func(og *Group) {
				var d time.Duration
				og.Duration(&d, "count", time.Nanosecond)
			},
			second: func(og *Group) {
				var d time.Duration
				og.DurationZero(&d, "count")
			},
		},
		"duplicate float64": {
			first: func(og *Group) {
				var f float64
				og.Float64(&f, "count", 1.0)
			},
			second: func(og *Group) {
				var f float64
				og.Float64Zero(&f, "count")
			},
		},
		"duplicate int": {
			first: func(og *Group) {
				var i int
				og.Int(&i, "count", 1)
			},
			second: func(og *Group) {
				var i int
				og.IntZero(&i, "count")
			},
		},
		"duplicate string": {
			first: func(og *Group) {
				var s string
				og.String(&s, "file", "first")
			},
			second: func(og *Group) {
				var s string
				og.StringZero(&s, "file")
			},
		},
		"duplicate uint": {
			first: func(og *Group) {
				var u uint
				og.Uint(&u, "count", 1)
			},
			second: func(og *Group) {
				var u uint
				og.UintZero(&u, "count")
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			og := NewGroup("test-optiongroup")
			tc.first(og)
			defer func() {
				if r := recover(); r == nil {
					t.Error("expected panic on duplicate registration")
				}
			}()
			tc.second(og)
		})
	}
}
