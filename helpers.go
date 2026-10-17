package opts

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Valid names must not be empty, begin with "-", or contain whitespace,
// control characters, quotes, backticks, backslashes, or equal signs.
// Whitespace and control characters are checked by Unicode class.
func isValidName(name string) bool {
	if name == "" || name[0] == '-' {
		return false
	}

	return !strings.ContainsFunc(name, isInvalidNameRune)
}

func isInvalidNameRune(r rune) bool {
	return unicode.IsSpace(r) || unicode.IsControl(r) || strings.ContainsRune("\"'`\\=", r)
}

func validateName(funcName, optName string) error {
	if !isValidName(optName) {
		return fmt.Errorf("opts: %s: invalid name: %q", funcName, optName)
	}

	return nil
}

// define registers an option of type T. funcName is the name of the calling
// method, for panic messages. If name is invalid or already defined, define
// panics without changing *p.
func define[T any](g *Group, funcName string, p *T, name string, defValue T, convert func(string) (T, error)) {
	if err := validateName(funcName, name); err != nil {
		panic(err)
	}

	if err := g.optAlreadyDefined(funcName, name); err != nil {
		panic(err)
	}

	// Only Bool defines options of type bool, and they take no value.
	_, isBool := any(p).(*bool)

	*p = defValue
	g.opts[name] = &opt{
		value:  &value[T]{ptr: p, convert: convert},
		isBool: isBool,
	}
}

func (g *Group) optAlreadyDefined(funcName, optName string) error {
	if _, exists := g.opts[optName]; exists {
		return fmt.Errorf("opts: %s: option %q already defined", funcName, optName)
	}

	return nil
}

// This function returns bare errors since the caller will give them context.
func numError(err error) error {
	ne, ok := errors.AsType[*strconv.NumError](err)
	if !ok {
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
