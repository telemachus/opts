package opts

import (
	"time"
)

// Duration defines a time.Duration option with the specified name and default
// value. The argument d points to a time.Duration variable that will store the
// value of the option. Duration will panic if name is not valid or repeats an
// existing option.
func (g *Group) Duration(d *time.Duration, name string, defValue time.Duration) {
	if err := validateName("Duration", name); err != nil {
		panic(err)
	}

	*d = defValue
	opt := &Opt{
		value: &value[time.Duration]{
			ptr:     d,
			convert: time.ParseDuration,
		},
		name:   name,
		isBool: false,
	}

	if err := g.optAlreadySet(name); err != nil {
		panic(err)
	}
	g.opts[name] = opt
}

// DurationZero is like Duration but with a default value of 0.
func (g *Group) DurationZero(d *time.Duration, name string) {
	g.Duration(d, name, 0)
}
