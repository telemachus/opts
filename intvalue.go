package opts

import (
	"strconv"
)

// Int defines an int option with the specified name and value. The argument
// i points to an int variable that will store the value of the option. Int
// will panic if name is not valid or repeats an existing option.
func (g *Group) Int(i *int, name string, defValue int) {
	if err := validateName("Int", name); err != nil {
		panic(err)
	}

	*i = defValue
	opt := &Opt{
		value: &value[int]{
			ptr:     i,
			convert: toInt,
		},
		name:   name,
		isBool: false,
	}

	if err := g.optAlreadySet(name); err != nil {
		panic(err)
	}
	g.opts[name] = opt
}

// IntZero defines an int option with the specified name and a value of 0. The
// argument i points to an int variable that will store the value of the
// option. IntZero will panic if name is not valid or repeats an existing
// option.
func (g *Group) IntZero(i *int, name string) {
	g.Int(i, name, 0)
}

func toInt(s string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		// The caller will give this error more context.
		return 0, numError(err)
	}
	return v, nil
}
