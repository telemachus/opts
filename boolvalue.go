package opts

import (
	"fmt"
)

type boolValue bool

func newBoolValue(p *bool) *boolValue {
	*p = false
	return (*boolValue)(p)
}

// Bool creates a bool option with the given name. Bool options always default
// to false. The argument p is a pointer to a bool variable that stores the
// value of the flag.
//
// Bool will panic if name is not a valid option name or if
// name repeats the name of an existing flag.
func (g *Group) Bool(b *bool, name string) {
	if err := validateName("Bool", name); err != nil {
		panic(err)
	}

	bv := newBoolValue(b)
	opt := &Opt{
		value:    bv,
		defValue: "false",
		name:     name,
		isBool:   true,
	}

	if err := g.optAlreadySet(name); err != nil {
		panic(err)
	}
	g.opts[name] = opt
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
