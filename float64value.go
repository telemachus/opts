package opts

import (
	"strconv"
)

// Float64 defines a float64 option with the specified name and default value.
// The argument f points to a float64 variable that will store the value of the
// option. Float64 will panic if name is not valid or repeats an existing
// option.
func (g *Group) Float64(f *float64, name string, defValue float64) {
	if err := validateName("Float64", name); err != nil {
		panic(err)
	}

	*f = defValue
	opt := &Opt{
		value: &value[float64]{
			ptr:     f,
			convert: toFloat64,
		},
		name:   name,
		isBool: false,
	}

	if err := g.optAlreadySet(name); err != nil {
		panic(err)
	}
	g.opts[name] = opt
}

// Float64Zero is like Float64 but with a default value of 0.0.
func (g *Group) Float64Zero(f *float64, name string) {
	g.Float64(f, name, 0.0)
}

func toFloat64(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, numError(err)
	}
	return v, nil
}
