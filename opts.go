package opts

import (
	"iter"
	"maps"
)

// Opt encapsulates a single option.
type Opt struct {
	value    setter
	defValue string
	name     string
	isBool   bool
}

// Setter is the interface that all options must satisfy. An option must be able
// to parse a string as a value and assign that value to a pointer variable of
// the appropriate type. If the string cannot be parsed as an appropriate
// value, the setter must return an error.
type setter interface {
	set(string) error
}

type value[T any] struct {
	ptr    *T
	parser func(string) (T, error)
}

func (v *value[T]) set(s string) error {
	val, err := v.parser(s)
	if err != nil {
		return err
	}

	*v.ptr = val

	return nil
}

// DefValue returns the default value of the option as a string.
func (o *Opt) DefValue() string {
	return o.defValue
}

// Name returns the command line name associated with the option.
func (o *Opt) Name() string {
	return o.name
}

// Group stores and manages a set of options.
type Group struct {
	opts   map[string]*Opt
	name   string
	args   []string // arguments remaining after option parsing
	parsed bool
}

// NewGroup returns a pointer to an option Group ready to use.
func NewGroup(name string) *Group {
	return &Group{
		name: name,
		opts: make(map[string]*Opt, 10),
	}
}

// Name returns the name of the group.
func (g *Group) Name() string {
	return g.name
}

// All iterates over the options in a group.
func (g *Group) All() iter.Seq[*Opt] {
	return maps.Values(g.opts)
}
