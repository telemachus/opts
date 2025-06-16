package opts

// String defines a string option with the specified name and default value.
// The argument s points to a string variable to hold the value of the option.
// String will panic if name is not valid or repeats an existing option.
func (g *Group) String(s *string, name, defValue string) {
	if err := validateName("String", name); err != nil {
		panic(err)
	}

	*s = defValue
	opt := &Opt{
		value: &genericValue[string]{
			target: s,
			parser: func(str string) (string, error) { return str, nil },
		},
		defValue: defValue,
		name:     name,
		isBool:   false,
	}

	if err := g.optAlreadySet(name); err != nil {
		panic(err)
	}
	g.opts[name] = opt
}

// StringZero defines a string option with the specified name and a default
// value of "". The argument s points to a string variable to hold the value of
// the option. StringZero will panic if name is not valid or repeats an
// existing option.
func (g *Group) StringZero(s *string, name string) {
	g.String(s, name, "")
}
