package opts

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// The following characters must not appear in the names of options.
//
// \x00				zero/null byte
// " " 				space
// \t\n\v\f\r\a\b\x1b		control characters
// \u0085			unicode next line (NEL)
// \u00A0			unicode non-breaking space (NBSP)
// "'				quotes
// `				backtick
// \				backslash
// =				equal
const junk = "\x00 \t\n\v\f\r\a\b\x1b\u0085\u00A0\"'`\\="

// Valid names must not be empty, begin with "-", or contain junk characters.
func isValidName(name string) bool {
	switch {
	case name == "":
		return false
	case name[0] == '-':
		return false
	case strings.ContainsAny(name, junk):
		return false
	default:
		return true
	}
}

func validateName(funcName, optName string) error {
	if !isValidName(optName) {
		return fmt.Errorf("opts: %s: invalid name: %s", funcName, optName)
	}

	return nil
}

func (g *Group) optAlreadySet(name string) error {
	if _, exists := g.opts[name]; exists {
		return fmt.Errorf("opts: --%s already set", name)
	}

	return nil
}

// This function returns bare errors since the caller will give them context.
func numError(err error) error {
	var ne *strconv.NumError
	if !errors.As(err, &ne) {
		return err
	}

	if errors.Is(ne.Err, strconv.ErrSyntax) {
		return strconv.ErrSyntax
	}

	if errors.Is(ne.Err, strconv.ErrRange) {
		return strconv.ErrRange
	}

	return err
}
