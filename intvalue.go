package opts

import (
	"fmt"
	"strconv"
)

// Int defines an int option with the specified name and value. The argument
// i points to an int variable to hold the value of the option. Int will panic
// if name is not valid or repeats an existing option.
func (g *Group) Int(i *int, name string, defValue int) {
	if err := validateName("Int", name); err != nil {
		panic(err)
	}

	*i = defValue
	opt := &Opt{
		value: &value[int]{
			ptr:    i,
			parser: parseInt,
		},
		defValue: strconv.Itoa(defValue),
		name:     name,
		isBool:   false,
	}

	if err := g.optAlreadySet(name); err != nil {
		panic(err)
	}
	g.opts[name] = opt
}

// IntZero defines an int option with the specified name and a value of 0. The
// argument i points to an int variable to hold the value of the option.
// IntZero will panic if name is not valid or repeats an existing option.
func (g *Group) IntZero(i *int, name string) {
	g.Int(i, name, 0)
}

func parseInt(s string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("parsing %q: %w", s, err)
	}
	return v, nil
}
