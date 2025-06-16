package opts

import (
	"errors"
	"fmt"
	"strconv"
)

// Float64 defines a float64 option with the specified name and default value.
// The argument f points to a float64 variable to hold the value of the option.
// Float64 will panic if name is not valid or repeats an existing option.
func (g *Group) Float64(f *float64, name string, defValue float64) {
	if err := validateName("Float64", name); err != nil {
		panic(err)
	}

	*f = defValue
	opt := &Opt{
		value: &value[float64]{
			ptr:    f,
			parser: parseFloat64,
		},
		defValue: strconv.FormatFloat(defValue, 'g', -1, 64),
		name:     name,
		isBool:   false,
	}

	if err := g.optAlreadySet(name); err != nil {
		panic(err)
	}
	g.opts[name] = opt
}

// Float64Zero defines a float64 option with the specified name and a default
// value of 0.0. The argument f points to a float64 variable to hold the value
// of the option. Float64Zero will panic if name is not valid or repeats an
// existing option.
func (g *Group) Float64Zero(f *float64, name string) {
	g.Float64(f, name, 0.0)
}

func parseFloat64(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, numError(err, s)
	}
	return v, nil
}

func numError(err error, s string) error {
	var ne *strconv.NumError
	if !errors.As(err, &ne) {
		return fmt.Errorf("%w", err)
	}

	if errors.Is(ne.Err, strconv.ErrSyntax) {
		return fmt.Errorf("parsing %q: %w", s, strconv.ErrSyntax)
	}

	if errors.Is(ne.Err, strconv.ErrRange) {
		return fmt.Errorf("parsing %q: %w", s, strconv.ErrRange)
	}

	return err
}
