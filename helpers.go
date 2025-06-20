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

func isValidName(name string) bool {
	isEmpty := name == ""
	initialDash := !isEmpty && name[0] == '-'
	hasJunk := strings.ContainsAny(name, junk)
	isValid := !isEmpty && !initialDash && !hasJunk

	return isValid
}

func validateName(funcName, optName string) error {
	if valid := isValidName(optName); !valid {
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
