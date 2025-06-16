package opts

import (
	"time"
)

// Duration defines a time.Duration option with the specified name and default
// value. The argument d points to a time.Duration variable to hold the value
// of the option. Duration will panic if name is not valid or repeats an
// existing option.
func (g *Group) Duration(d *time.Duration, name string, defValue time.Duration) {
	if err := validateName("Duration", name); err != nil {
		panic(err)
	}

	*d = defValue
	opt := &Opt{
		value: &genericValue[time.Duration]{
			target: d,
			parser: time.ParseDuration,
		},
		defValue: defValue.Round(time.Second).String(),
		name:     name,
		isBool:   false,
	}

	if err := g.optAlreadySet(name); err != nil {
		panic(err)
	}
	g.opts[name] = opt
}

// DurationZero defines a time.Duration option with the specified name and
// a default value of 0. The argument d points to a time.Duration variable to
// hold the value of the option. DurationZero will panic if name is not valid
// or repeats an existing option.
func (g *Group) DurationZero(d *time.Duration, name string) {
	g.Duration(d, name, 0)
}
