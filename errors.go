package opts

import (
	"errors"
	"fmt"
)

// ErrAlreadyParsed signals an attempt to parse a [Group] that has already been
// successfully parsed.
var ErrAlreadyParsed = errors.New("opts: option group already parsed")

// UnexpectedArgsError signals that there are args left after parsing. Only
// [Parse] will return this error. Use [ParseKnown] for relaxed parsing.
type UnexpectedArgsError struct {
	Args []string
}

func (e *UnexpectedArgsError) Error() string {
	var s string
	if len(e.Args) > 1 {
		s = "s"
	}

	return fmt.Sprintf("opts: unexpected argument%s after parsing: %v", s, e.Args)
}

// MissingValueError signals that an option is missing a required value.
type MissingValueError struct {
	Option string
}

func (e *MissingValueError) Error() string {
	return fmt.Sprintf("opts: missing value for --%s", e.Option)
}

// InvalidValueError signals that an option's value cannot be converted into
// the option's type. Since InvalidValueError wraps the original conversion
// error, users can access the undedited original as InvalidValueError.Err.
type InvalidValueError struct {
	Err    error
	Option string
	Value  string
}

func (e *InvalidValueError) Error() string {
	return fmt.Sprintf("opts: invalid value %q for --%s: %v", e.Value, e.Option, e.Err)
}

func (e *InvalidValueError) Unwrap() error {
	return e.Err
}

// UnknownOptionError signals that an option was not registered with the Group.
type UnknownOptionError struct {
	Option string
}

func (e *UnknownOptionError) Error() string {
	return fmt.Sprintf("opts: unknown option --%s", e.Option)
}
