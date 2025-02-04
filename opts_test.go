package opts

import (
	"testing"
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

func TestBoolSet(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input string
		want  bool
	}{
		"true":  {input: "true", want: true},
		"false": {input: "false", want: false},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var b boolValue
			if err := b.set(tc.input); err != nil {
				t.Fatalf("b.set(%q) unexpected error: %q", tc.input, err)
			}

			if bool(b) != tc.want {
				t.Errorf("after b.set(%q), b == %v; want %v", tc.input, b, tc.want)
			}
		})
	}
}

func TestBoolValueSetError(t *testing.T) {
	t.Parallel()

	badInputs := []string{
		"",
		"yeah",
		"nope",
		"2",
		"-1",
	}

	for _, input := range badInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()

			var b boolValue
			if err := b.set(input); err == nil {
				t.Errorf("b.set(%q) err == nil; want error", input)
			}
		})
	}
}

//
// func TestDurationValueSet(t *testing.T) {
// 	t.Parallel()
//
// 	testCases := map[string]struct {
// 		input string
// 		want  time.Duration
// 	}{
// 		"zero":           {input: "0s", want: 0},
// 		"seconds":        {input: "10s", want: 10 * time.Second},
// 		"minutes":        {input: "5m", want: 5 * time.Minute},
// 		"hours":          {input: "2h", want: 2 * time.Hour},
// 		"negative":       {input: "-30s", want: -30 * time.Second},
// 		"mixed":          {input: "2h30m", want: 2*time.Hour + 30*time.Minute},
// 		"milliseconds":   {input: "1500ms", want: 1500 * time.Millisecond},
// 		"microseconds":   {input: "1500µs", want: 1500 * time.Microsecond},
// 		"multiple units": {input: "1h2h", want: 3 * time.Hour},
// 		"nanoseconds":    {input: "1500ns", want: 1500 * time.Nanosecond},
// 	}
//
// 	for msg, tc := range testCases {
// 		t.Run(msg, func(t *testing.T) {
// 			t.Parallel()
//
// 			var d durationValue
// 			if err := d.set(tc.input); err != nil {
// 				t.Fatalf("d.set(%q) unexpected error: %q", tc.input, err)
// 			}
//
// 			if time.Duration(d) != tc.want {
// 				t.Errorf("after d.set(%q), d == %v; want %v", tc.input, d, tc.want)
// 			}
// 		})
// 	}
// }
//
// func TestDurationValueSetError(t *testing.T) {
// 	t.Parallel()
//
// 	badInputs := map[string]string{
// 		"empty string":     "",
// 		"just number":      "42",
// 		"invalid unit":     "42x",
// 		"invalid format":   "1.5.h",
// 		"letters":          "abc",
// 		"mixed letters":    "10xyz",
// 		"missing number":   "h",
// 		"spaces":           "1 h",
// 		"invalid negative": "-h",
// 	}
//
// 	for msg, input := range badInputs {
// 		t.Run(msg, func(t *testing.T) {
// 			t.Parallel()
//
// 			var d durationValue
// 			if err := d.set(input); err == nil {
// 				t.Errorf("d.set(%q) err == nil; want error", input)
// 			}
// 		})
// 	}
// }
//
// func TestFloat64ValueSet(t *testing.T) {
// 	t.Parallel()
//
// 	testCases := map[string]struct {
// 		input string
// 		want  float64
// 	}{
// 		"zero":           {input: "0", want: 0.0},
// 		"positive":       {input: "42", want: 42.0},
// 		"negative":       {input: "-123.123", want: -123.123},
// 		"leading zeros":  {input: "000123", want: 123.0},
// 		"plus sign":      {input: "+456", want: 456.0},
// 		"large positive": {input: "999999999", want: 999999999.0},
// 		"large negative": {input: "-999999999", want: -999999999.0},
// 	}
//
// 	for msg, tc := range testCases {
// 		t.Run(msg, func(t *testing.T) {
// 			t.Parallel()
//
// 			var f float64Value
// 			if err := f.set(tc.input); err != nil {
// 				t.Fatalf("f.set(%q) unexpected error: %q", tc.input, err)
// 			}
//
// 			if float64(f) != tc.want {
// 				t.Errorf("after f.set(%q), f == %v; want %v", tc.input, f, tc.want)
// 			}
// 		})
// 	}
// }
//
// func TestFloat64ValueSetError(t *testing.T) {
// 	t.Parallel()
//
// 	badInputs := map[string]string{
// 		"empty string":                      "",
// 		"letters: abc":                      "abc",
// 		"multiple decimal points: 12.34.24": "12.34.24",
// 		"invalid scientific notation: e10":  "e10",
// 		"invalid scientific notation: 1e":   "1e",
// 		"numbers and letters: 123abc":       "123abc",
// 		"overflow: 1e+1000":                 "1e+1000",
// 	}
//
// 	for msg, input := range badInputs {
// 		t.Run(msg, func(t *testing.T) {
// 			t.Parallel()
//
// 			var f float64Value
// 			if err := f.set(input); err == nil {
// 				t.Errorf("f.set(%q) err == nil; want error", input)
// 			}
// 		})
// 	}
// }
//
// func TestIntValueSet(t *testing.T) {
// 	t.Parallel()
//
// 	testCases := map[string]struct {
// 		input string
// 		want  int
// 	}{
// 		"zero":           {input: "0", want: 0},
// 		"positive":       {input: "42", want: 42},
// 		"negative":       {input: "-123", want: -123},
// 		"max int32":      {input: "2147483647", want: math.MaxInt32},
// 		"min int32":      {input: "-2147483648", want: math.MinInt32},
// 		"leading zeros":  {input: "000123", want: 123},
// 		"plus sign":      {input: "+456", want: 456},
// 		"large positive": {input: "999999999", want: 999999999},
// 		"large negative": {input: "-999999999", want: -999999999},
// 	}
//
// 	for msg, tc := range testCases {
// 		t.Run(msg, func(t *testing.T) {
// 			t.Parallel()
//
// 			var i intValue
// 			if err := i.set(tc.input); err != nil {
// 				t.Fatalf("i.set(%q) unexpected error: %q", tc.input, err)
// 			}
//
// 			if int(i) != tc.want {
// 				t.Errorf("after i.set(%q), i == %v; want %v", tc.input, i, tc.want)
// 			}
// 		})
// 	}
// }
//
// func TestIntValueSetError(t *testing.T) {
// 	t.Parallel()
//
// 	badInputs := map[string]string{
// 		"empty string":                    "",
// 		"letters: abc":                    "abc",
// 		"float: 12.34":                    "12.34",
// 		"scientific notation: 1e6":        "1e6",
// 		"hex: 0x123":                      "0x123",
// 		"underscore: 1_000":               "1_000",
// 		"word: nine":                      "nine",
// 		"NaN":                             "NaN",
// 		"Inf":                             "Inf",
// 		"-Inf":                            "-Inf",
// 		"overflow: 9223372036854775808":   "9223372036854775808",
// 		"underflow: -9223372036854775809": "-9223372036854775809",
// 	}
//
// 	for msg, input := range badInputs {
// 		t.Run(msg, func(t *testing.T) {
// 			t.Parallel()
//
// 			var i intValue
// 			if err := i.set(input); err == nil {
// 				t.Errorf("i.set(%q) err == nil; want error", input)
// 			}
// 		})
// 	}
// }
//
// func TestStringValueSet(t *testing.T) {
// 	t.Parallel()
//
// 	testCases := map[string]struct {
// 		input string
// 		want  string
// 	}{
// 		"empty string":  {input: "", want: ""},
// 		"simple string": {input: "test", want: "test"},
// 		"with spaces":   {input: "hello world", want: "hello world"},
// 		"unicode":       {input: "テスト", want: "テスト"},
// 		"special chars": {input: "!@#$%^&*()", want: "!@#$%^&*()"},
// 		"very long":     {input: strings.Repeat("x", 1000), want: strings.Repeat("x", 1000)},
// 		"with newlines": {input: "line1\nline2", want: "line1\nline2"},
// 		"with tabs":     {input: "col1\tcol2", want: "col1\tcol2"},
// 	}
//
// 	for msg, tc := range testCases {
// 		t.Run(msg, func(t *testing.T) {
// 			t.Parallel()
//
// 			var s stringValue
// 			if err := s.set(tc.input); err != nil {
// 				t.Fatalf("s.set(%q) unexpected error: %q", tc.input, err)
// 			}
//
// 			if string(s) != tc.want {
// 				t.Errorf("after s.set(%q), s == %q; want %q", tc.input, s, tc.want)
// 			}
// 		})
// 	}
// }
//
// func TestUintValueSet(t *testing.T) {
// 	t.Parallel()
//
// 	testCases := map[string]struct {
// 		input string
// 		want  uint
// 	}{
// 		"zero":         {input: "0", want: 0},
// 		"positive":     {input: "42", want: 42},
// 		"hex value":    {input: "0xff", want: 255},
// 		"octal value":  {input: "0644", want: 420},
// 		"large number": {input: "999999999", want: 999999999},
// 	}
//
// 	for msg, tc := range testCases {
// 		t.Run(msg, func(t *testing.T) {
// 			t.Parallel()
//
// 			var u uintValue
// 			if err := u.set(tc.input); err != nil {
// 				t.Fatalf("u.set(%q) unexpected error: %q", tc.input, err)
// 			}
//
// 			if uint(u) != tc.want {
// 				t.Errorf("after u.set(%q), u == %v; want %v", tc.input, u, tc.want)
// 			}
// 		})
// 	}
// }
//
// func TestUintValueSetError(t *testing.T) {
// 	t.Parallel()
//
// 	badInputs := map[string]string{
// 		"empty string": "",
// 		"letters":      "abc",
// 		"float":        "12.34",
// 		"negative":     "-42",
// 		"mixed":        "123abc",
// 		"overflow":     "18446744073709551616", // uint64 max + 1
// 		"symbols":      "!@#",
// 		"spaces":       "42 43",
// 	}
//
// 	for msg, input := range badInputs {
// 		t.Run(msg, func(t *testing.T) {
// 			t.Parallel()
//
// 			var u uintValue
// 			if err := u.set(input); err == nil {
// 				t.Errorf("u.set(%q) err == nil; want error", input)
// 			}
// 		})
// 	}
// }
//
// func TestFlagRegistrationBothValid(t *testing.T) {
// 	t.Parallel()
//
// 	var (
// 		short    = 'v'
// 		long     = "verbose"
// 		defValue = true
// 	)
//
// 	fs := NewFlagSet("test-flagset")
// 	var got bool
// 	fs.BoolVar(&got, short, long, defValue)
//
// 	if flag := fs.shorts[short]; flag == nil {
// 		t.Errorf("short flag -%c not registered", short)
// 	}
//
// 	if flag := fs.longs[long]; flag == nil {
// 		t.Errorf("long flag --%s not registered", long)
// 	}
//
// 	if fs.shorts[short] != fs.longs[long] {
// 		t.Error("short and long flags point to different Flag instances")
// 	}
// }
//
// func TestFlagRegistrationShortValid(t *testing.T) {
// 	t.Parallel()
//
// 	var (
// 		short    = 'v'
// 		long     = ""
// 		defValue = true
// 	)
//
// 	fs := NewFlagSet("test-flagset")
// 	var got bool
// 	fs.BoolVar(&got, short, long, defValue)
//
// 	if flag := fs.shorts[short]; flag == nil {
// 		t.Errorf("short flag -%c not registered", short)
// 	}
//
// 	if _, exists := fs.longs[long]; exists {
// 		t.Errorf("invalid long flag --%s was registered", long)
// 	}
// }
//
// func TestFlagRegistrationLongValid(t *testing.T) {
// 	t.Parallel()
// 	var (
// 		short    rune
// 		long     = "verbose"
// 		defValue = true
// 	)
//
// 	fs := NewFlagSet("test-flagset")
// 	var got bool
// 	fs.BoolVar(&got, short, long, defValue)
//
// 	if flag := fs.longs[long]; flag == nil {
// 		t.Errorf("long flag --%s not registered", long)
// 	}
//
// 	if _, exists := fs.shorts[short]; exists {
// 		t.Errorf("invalid short flag -%c was registered", short)
// 	}
// }
//
// func TestDuplicateFlagRegistration(t *testing.T) {
// 	t.Parallel()
//
// 	testCases := map[string]struct {
// 		first  func(*FlagSet)
// 		second func(*FlagSet)
// 	}{
// 		"duplicate short bool": {
// 			first: func(fs *FlagSet) {
// 				var b bool
// 				fs.BoolVar(&b, 'v', "", false)
// 			},
// 			second: func(fs *FlagSet) {
// 				var b bool
// 				fs.BoolVar(&b, 'v', "", true)
// 			},
// 		},
// 		"duplicate long bool": {
// 			first: func(fs *FlagSet) {
// 				var b bool
// 				fs.BoolVar(&b, 0, "verbose", false)
// 			},
// 			second: func(fs *FlagSet) {
// 				var b bool
// 				fs.BoolVar(&b, 0, "verbose", true)
// 			},
// 		},
// 		"duplicate short int": {
// 			first: func(fs *FlagSet) {
// 				var i int
// 				fs.IntVar(&i, 'n', "", 1)
// 			},
// 			second: func(fs *FlagSet) {
// 				var i int
// 				fs.IntVar(&i, 'n', "", 2)
// 			},
// 		},
// 		"duplicate long int": {
// 			first: func(fs *FlagSet) {
// 				var i int
// 				fs.IntVar(&i, 0, "count", 1)
// 			},
// 			second: func(fs *FlagSet) {
// 				var i int
// 				fs.IntVar(&i, 0, "count", 2)
// 			},
// 		},
// 		"duplicate short string": {
// 			first: func(fs *FlagSet) {
// 				var s string
// 				fs.StringVar(&s, 'f', "", "first")
// 			},
// 			second: func(fs *FlagSet) {
// 				var s string
// 				fs.StringVar(&s, 'f', "", "second")
// 			},
// 		},
// 		"duplicate long string": {
// 			first: func(fs *FlagSet) {
// 				var s string
// 				fs.StringVar(&s, 0, "file", "first")
// 			},
// 			second: func(fs *FlagSet) {
// 				var s string
// 				fs.StringVar(&s, 0, "file", "second")
// 			},
// 		},
// 	}
//
// 	for name, tc := range testCases {
// 		t.Run(name, func(t *testing.T) {
// 			t.Parallel()
// 			fs := NewFlagSet("test-flagset")
//
// 			tc.first(fs)
// 			defer func() {
// 				if r := recover(); r == nil {
// 					t.Error("expected panic on duplicate registration")
// 				}
// 			}()
// 			tc.second(fs)
// 		})
// 	}
// }
//
// func TestBool(t *testing.T) {
// 	t.Parallel()
//
// 	fs := NewFlagSet("test")
// 	b := fs.Bool('v', "verbose", false)
// 	if b == nil {
// 		t.Fatalf("fs.Bool('v', %q, false) returned nil", "verbose")
// 	}
// 	if *b != false {
// 		t.Errorf("*b == %t; want false", *b)
// 	}
// }
//
// func TestBoolShort(t *testing.T) {
// 	t.Parallel()
//
// 	fs := NewFlagSet("test")
// 	b := fs.BoolShort('v', true)
// 	if b == nil {
// 		t.Fatal("fs.Bool('v', true) returned nil")
// 	}
// 	if *b != true {
// 		t.Errorf("*b == %t; want true", *b)
// 	}
//
// 	if len(fs.shorts) != 1 {
// 		t.Error("short flag was not registered")
// 	}
// 	if len(fs.longs) != 0 {
// 		t.Error("long flag was registered")
// 	}
// }
//
// func TestBoolVarShort(t *testing.T) {
// 	t.Parallel()
//
// 	fs := NewFlagSet("test")
// 	var got bool
// 	fs.BoolVarShort(&got, 'v', false)
//
// 	err := fs.Parse([]string{"-v"})
// 	if err != nil {
// 		t.Fatalf("fs.Parse(%q) failed: %v", "-v", err)
// 	}
// 	if !got {
// 		t.Error("flag not set to true")
// 	}
//
// 	if len(fs.shorts) != 1 {
// 		t.Error("short flag was not registered")
// 	}
// 	if len(fs.longs) != 0 {
// 		t.Error("long flag was registered")
// 	}
// }
//
// func TestBoolVarLong(t *testing.T) {
// 	t.Parallel()
//
// 	fs := NewFlagSet("test")
// 	var got bool
// 	fs.BoolVarLong(&got, "verbose", false)
//
// 	err := fs.Parse([]string{"--verbose"})
// 	if err != nil {
// 		t.Fatalf("fs.Parse(%q) failed: %v", "--verbose", err)
// 	}
// 	if !got {
// 		t.Error("flag not set to true")
// 	}
//
// 	if len(fs.longs) != 1 {
// 		t.Error("long flag was not registered")
// 	}
// 	if len(fs.shorts) != 0 {
// 		t.Error("short flag was registered")
// 	}
// }
