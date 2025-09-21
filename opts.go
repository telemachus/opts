package opts

// An opt stores a single option.
type opt struct {
	value  setter
	name   string
	isBool bool
}

// Options implement the setter interface, parsing a given string and assigning
// its value to a pointer of the option's type or returning an error if parsing
// fails.
type setter interface {
	set(string) error
}

type value[T any] struct {
	ptr     *T
	convert func(string) (T, error)
}

func (v *value[T]) set(s string) error {
	val, err := v.convert(s)
	if err != nil {
		return err
	}

	*v.ptr = val

	return nil
}

// Group stores and manages a set of options.
type Group struct {
	opts   map[string]*opt
	name   string
	args   []string
	parsed bool
}

// NewGroup returns a pointer to an option Group ready to use.
func NewGroup(name string) *Group {
	return &Group{
		name: name,
		opts: make(map[string]*opt, 10),
	}
}

// Name returns the name of the group.
func (g *Group) Name() string {
	return g.name
}
