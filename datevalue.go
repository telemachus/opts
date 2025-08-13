package opts

import (
	"cloud.google.com/go/civil"
)

// Date defines a civil.Date option with the specified name and default value.
// The argument d points to a civil.Date variable that will store the value of
// the option. Date will panic if name is not valid or repeats an existing
// option. On the command line, users must use a string in RFC 3339 full-date
// format (e.g., "2025-12-31" or "2024-02-29"). See [civil.ParseDate] for
// details.
func (g *Group) Date(d *civil.Date, name string, defValue civil.Date) {
	if err := validateName("Date", name); err != nil {
		panic(err)
	}

	*d = defValue
	opt := &Opt{
		value: &value[civil.Date]{
			ptr:     d,
			convert: civil.ParseDate,
		},
		name:   name,
		isBool: false,
	}

	if err := g.optAlreadySet(name); err != nil {
		panic(err)
	}
	g.opts[name] = opt
}

// DateZero is like Date but with a zero default value.
func (g *Group) DateZero(d *civil.Date, name string) {
	g.Date(d, name, civil.Date{})
}
