package opts

import (
	"fmt"
	"strconv"
)

type boolValue bool

func newBoolValue(val bool, b *bool) *boolValue {
	*b = val
	return (*boolValue)(b)
}

// Bool creates a new boolean option and binds its default value to b. Bool
// will panic if name is not a valid option name or if name repeats the name of
// an existing option.
//
// NOTE: this method is not exported. Only BoolZero is exported because opts
// always treats boolean options as switches. They default to false, and if
// a user passes the option on the command line, the boolean value becomes
// true. Boolean options do not accept arguments on the command line.
func (g *Group) bool(b *bool, name string, defValue bool) {
	if err := validateName("Bool", name); err != nil {
		panic(err)
	}

	bv := newBoolValue(defValue, b)
	opt := &Opt{
		value:    bv,
		defValue: strconv.FormatBool(defValue),
		name:     name,
		isBool:   true,
	}

	if err := g.optAlreadySet(name); err != nil {
		panic(err)
	}
	g.opts[name] = opt
}

// BoolZero creates a new boolean option that defaults to false.
func (g *Group) BoolZero(b *bool, name string) {
	g.bool(b, name, false)
}

func parseBool(s string) (bool, error) {
	switch s {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("cannot parse %q", s)
	}
}

// TODO: restrict valid boolean values to "true" and "false"?
//
// Set assigns s to an boolValue and returns an error if s cannot be parsed as
// a boolean. Valid boolean values are 1, 0, t, f, T, F, true, false, TRUE,
// FALSE, True, False.
func (b *boolValue) set(s string) error {
	v, err := parseBool(s)
	if err != nil {
		return err
	}

	*b = boolValue(v)

	return nil
}
