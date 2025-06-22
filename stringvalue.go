package opts

// String defines a string option with the specified name and default value.
// The argument s points to a string variable that will store the value of the
// option. String will panic if name is not valid or repeats an existing
// option.
func (g *Group) String(s *string, name, defValue string) {
	if err := validateName("String", name); err != nil {
		panic(err)
	}

	*s = defValue
	opt := &Opt{
		value: &value[string]{
			ptr:     s,
			convert: toString,
		},
		name:   name,
		isBool: false,
	}

	if err := g.optAlreadySet(name); err != nil {
		panic(err)
	}
	g.opts[name] = opt
}

// StringZero is like String but with a default value of "".
func (g *Group) StringZero(s *string, name string) {
	g.String(s, name, "")
}

func toString(s string) (string, error) {
	return s, nil
}
