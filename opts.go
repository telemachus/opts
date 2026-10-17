package opts

// An opt stores a single option.
type opt struct {
	value  preparer
	isBool bool
}

// Options implement the preparer interface, parsing a given string and
// returning a function that assigns the parsed value to a pointer of the
// option's type, or returning an error if parsing fails. Nothing is assigned
// until the caller runs the returned function.
type preparer interface {
	prepare(s string) (func(), error)
}

type value[T any] struct {
	ptr     *T
	convert func(string) (T, error)
}

func (v *value[T]) prepare(s string) (func(), error) {
	val, err := v.convert(s)
	if err != nil {
		return nil, err
	}

	return func() { *v.ptr = val }, nil
}

// Group stores and manages a set of options.
type Group struct {
	opts   map[string]*opt
	name   string
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
