package opts

import (
	"fmt"
)

// Bool defines a bool option with the specified name and a default value of
// false. The argument b points to a bool variable that will store the value.
// Bool will panic if name is not valid or repeats an existing option.
func (g *Group) Bool(b *bool, name string) {
	if err := validateName("Bool", name); err != nil {
		panic(err)
	}

	*b = false
	opt := &opt{
		value: &value[bool]{
			ptr:     b,
			convert: toBool,
		},
		name:   name,
		isBool: true,
	}

	if err := g.optAlreadySet(name); err != nil {
		panic(err)
	}
	g.opts[name] = opt
}

func toBool(s string) (bool, error) {
	switch s {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		// Omit "opts: " since the caller will provide context.
		return false, fmt.Errorf("bool value must be %q or %q", "true", "false")
	}
}
