package opts

import (
	"fmt"
)

// Bool defines a bool option with the specified name and a value of false. The
// argument b points to a bool variable to hold the value of the option. Bool
// will panic if name is not valid or repeats an existing option.
func (g *Group) Bool(b *bool, name string) {
	if err := validateName("Bool", name); err != nil {
		panic(err)
	}

	*b = false
	opt := &Opt{
		value: &value[bool]{
			ptr:    b,
			parser: parseBool,
		},
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
