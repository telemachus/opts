package opts_test

import (
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/telemachus/opts"
)

func TestNewGroup(t *testing.T) {
	t.Parallel()

	name := "test-opts"
	og := opts.NewGroup(name)
	if og == nil {
		t.Fatalf("NewGroup(%q) returned nil", name)
	}

	if og.Name() != name {
		t.Errorf("og.Name(%q) == %q; want %q", name, og.Name(), name)
	}
}

func TestOptRegistrationInvalid(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		name string
	}{
		"empty name":         {name: ""},
		"whitespace":         {name: " hello"},
		"equal sign":         {name: "hello=world"},
		"tab":                {name: "hello\tworld"},
		"newline":            {name: "hello\nworld"},
		"null byte":          {name: "hello\u0000world"},
		"null rune":          {name: "hello\x00world"},
		"initial dash":       {name: "-hello"},
		"backslash":          {name: `hello\world`},
		"single quote":       {name: "hello'world"},
		"double quote":       {name: `hello"world`},
		"backtick":           {name: "hello`world"},
		"unicode whitespace": {name: "hello\u00A0world"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			og := opts.NewGroup("test-optiongroup")
			defer func() {
				if r := recover(); r == nil {
					t.Error("expected panic on invalid name")
				}
			}()
			var got bool
			og.Bool(&got, tc.name)
		})
	}
}

func TestDuplicateOptRegistration(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		first  func(*opts.Group)
		second func(*opts.Group)
	}{
		"duplicate bool": {
			first: func(og *opts.Group) {
				var b bool
				og.Bool(&b, "verbose")
			},
			second: func(og *opts.Group) {
				var b bool
				og.Bool(&b, "verbose")
			},
		},
		"duplicate date": {
			first: func(og *opts.Group) {
				var d civil.Date
				og.Date(&d, "birthday", civil.Date{Year: 1972, Month: 6, Day: 23})
			},
			second: func(og *opts.Group) {
				var d civil.Date
				og.Date(&d, "birthday", civil.Date{Year: 1972, Month: 6, Day: 23})
			},
		},
		"duplicate duration": {
			first: func(og *opts.Group) {
				var d time.Duration
				og.Duration(&d, "count", time.Nanosecond)
			},
			second: func(og *opts.Group) {
				var d time.Duration
				og.DurationZero(&d, "count")
			},
		},
		"duplicate float64": {
			first: func(og *opts.Group) {
				var f float64
				og.Float64(&f, "count", 1.0)
			},
			second: func(og *opts.Group) {
				var f float64
				og.Float64Zero(&f, "count")
			},
		},
		"duplicate int": {
			first: func(og *opts.Group) {
				var i int
				og.Int(&i, "count", 1)
			},
			second: func(og *opts.Group) {
				var i int
				og.IntZero(&i, "count")
			},
		},
		"duplicate string": {
			first: func(og *opts.Group) {
				var s string
				og.String(&s, "file", "first")
			},
			second: func(og *opts.Group) {
				var s string
				og.StringZero(&s, "file")
			},
		},
		"duplicate uint": {
			first: func(og *opts.Group) {
				var u uint
				og.Uint(&u, "count", 1)
			},
			second: func(og *opts.Group) {
				var u uint
				og.UintZero(&u, "count")
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			og := opts.NewGroup("test-optiongroup")
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
