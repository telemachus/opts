package opts

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrAlreadyParsed signals an attempt to parse a [Group] that has already been
// successfully parsed.
var ErrAlreadyParsed = errors.New("already parsed")

// ErrBooleanWithValue signals an unsupported use of "option=value" with
// a boolean.
var ErrBooleanWithValue = errors.New("value given for boolean option")

// ErrMissingValue signals that an option is missing a required value.
var ErrMissingValue = errors.New("missing value")

// ErrUnknownOption signals that an option was not registered with the [Group].
var ErrUnknownOption = errors.New("unknown option")

// UnexpectedArgumentsError signals that there are one or more arguments left
// after parsing. Only [Group.ParseAll] will return this error. Use
// [Group.Parse] to accept positional arguments after the options.
type UnexpectedArgumentsError struct {
	Args []string
}

func (e *UnexpectedArgumentsError) Error() string {
	var s string
	if len(e.Args) > 1 {
		s = "s"
	}

	return fmt.Sprintf("opts: unexpected argument%s: %s", s, quotedArgs(e.Args))
}

// We can rely on Args having at least one item since ParseAll does not return
// UnexpectedArgumentsError if there are no leftover arguments.
func quotedArgs(items []string) string {
	var b strings.Builder
	b.WriteString(strconv.Quote(items[0]))
	for _, str := range items[1:] {
		b.WriteString(", ")
		b.WriteString(strconv.Quote(str))
	}

	return b.String()
}

// InvalidValueError signals that an option's value cannot be converted into
// the option's type. Err holds the reason that conversion failed. For
// numeric options, Err is [strconv.ErrSyntax] or [strconv.ErrRange]. Option
// is the option's name as defined, without dashes.
type InvalidValueError struct {
	Err    error
	Option string
	Value  string

	// optPrefix records the single or double dash before an option.
	// Error uses it to display dashes as the user typed them.
	optPrefix string
}

func (e *InvalidValueError) Error() string {
	return fmt.Sprintf("opts: invalid value %q for %q: %v", e.Value, e.optPrefix+e.Option, e.Err)
}

func (e *InvalidValueError) Unwrap() error {
	return e.Err
}
