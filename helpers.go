package opts

import (
	"fmt"
	"strings"
)

var junk = strings.Join([]string{
	string(rune(0x00)),   // zero rune
	string([]byte{0x00}), // zero/NUL byte
	` `,                  // space
	"\t\n\v\f\r",         // control whitespace
	string(rune(0x85)),   // unicode whitespace
	string(rune(0xA0)),   // unicode whitespace
	`"'`,                 // quotes
	"`",                  // backtick
	`\`,                  // backslash
}, "")

func isValidName(name string) bool {
	isEmpty := name == ""
	hasJunk := strings.ContainsAny(name, junk)
	isValid := !isEmpty && !hasJunk

	return isValid
}

func validateName(funcName, optName string) error {
	if valid := isValidName(optName); !valid {
		return fmt.Errorf("opt: %s: invalid name: %s", funcName, optName)
	}

	return nil
}

func (g *Group) optAlreadySet(name string) error {
	if _, exists := g.opts[name]; exists {
		// TODO: quote the option?
		return fmt.Errorf("opt: --%s already set", name)
	}

	return nil
}
