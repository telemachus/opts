package opts

type genericValue[T any] struct {
	target *T
	parser func(string) (T, error)
}

func (g *genericValue[T]) set(s string) error {
	val, err := g.parser(s)
	if err != nil {
		return err
	}

	*g.target = val

	return nil
}
