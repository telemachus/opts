package opts

import (
	"strconv"
)

// Uint defines a uint option with the specified name and default value. The
// argument u points to a uint variable to hold the value of the option. Uint
// will panic if name is not valid or repeats an existing option.
func (g *Group) Uint(u *uint, name string, defValue uint) {
	if err := validateName("Uint", name); err != nil {
		panic(err)
	}

	*u = defValue
	opt := &Opt{
		value: &genericValue[uint]{
			target: u,
			parser: parseUint,
		},
		defValue: strconv.FormatUint(uint64(defValue), 10),
		name:     name,
		isBool:   false,
	}

	if err := g.optAlreadySet(name); err != nil {
		panic(err)
	}
	g.opts[name] = opt
}

// UintZero defines a uint option with the specified name and default value.
// The argument u points to a uint variable to hold the value of the option.
// UintZero will panic if name is not valid or repeats an existing option.
func (g *Group) UintZero(u *uint, name string) {
	g.Uint(u, name, 0)
}

func parseUint(s string) (uint, error) {
	v, err := strconv.ParseUint(s, 0, strconv.IntSize)
	if err != nil {
		return 0, numError(err, s)
	}
	return uint(v), nil
}
