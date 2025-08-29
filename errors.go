package opts

import (
	"errors"
	"fmt"
)

// ErrAlreadyParsed is returned when [Parse] or [ParseKnown] is called on
// a [Group] that has already been successfully parsed.
var ErrAlreadyParsed = errors.New("opts: option group already parsed")

// ErrUnexpectedArgs is returned when [Parse] is called and there are remaining
// args left that have not be registered as options.
var ErrUnexpectedArgs = errors.New("opts: unexpected arguments remain after parsing")

// ErrMissingValue is returned when an option that requires a value is not
// followed by one.
var ErrMissingValue = errors.New("option requires a value")

// InvalidValueError is returned when an option string cannot be converted into
// the correct type. It wraps the underlying conversion error.
type InvalidValueError struct {
	Err    error
	Option string
	Value  string
}

func (e *InvalidValueError) Error() string {
	return fmt.Sprintf("opts: invalid value %q for option %s: %v", e.Value, e.Option, e.Err)
}

func (e *InvalidValueError) Unwrap() error {
	return e.Err
}
