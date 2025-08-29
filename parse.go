package opts

import (
	"fmt"
	"strings"
)

// Parse scans args and sets option values defined in the [Group].
//
// Parse should be called after all options are defined and before any option
// values are used. Parse performs strict parsing. If any positional arguments
// remain after parsing stops, then Parse will return [ErrUnexpectedArgs]. If
// Parse returns without error, the Group is considered parsed and subsequent
// calls to Parse will return [ErrAlreadyParsed].
//
// Parsing stops at the first non-option argument. Any remaining arguments can
// be accessed via the returned slice of strings. Both '-' and '--' are
// considered non-option arguments, and both stop further parsing. However,
// the returned slice will not contain '--', but it will contain '-'. (By
// convention, many programs treat '-' as stdin, but that is up to the calling
// program to decide and handle.)
//
// Parse performs strict parsing. If any positional arguments remain after
// parsing stops, then Parse will return [ErrUnexpectedArgs].
//
// If Parse encounters an unknown option, an option without a value, or a value
// that cannot be parsed as its type, it returns an error and the Group remains
// unparsed. The caller may retry with different arguments.
//
// The slice passed to Parse should not include the program name. If using
// `os.Args` directly, the caller should pass `os.Args[1:]`.
//
// Use Parse if your program should not accept leftover arguments.
func (g *Group) Parse(args []string) ([]string, error) {
	if g.parsed {
		return []string{}, ErrAlreadyParsed
	}

	err := g.parse(args)
	if err != nil {
		return []string{}, err
	}

	if len(g.args) > 0 {
		return g.args, ErrUnexpectedArgs
	}

	g.parsed = true

	return g.args, nil
}

// ParseKnown is like Parse in all ways but one: it allows arguments to remain
// after parsing stops and thus never returns ErrUnexpectedArgs.
//
// Use ParseKnown if your program expects leftover arguments.
func (g *Group) ParseKnown(args []string) ([]string, error) {
	if g.parsed {
		return []string{}, ErrAlreadyParsed
	}

	err := g.parse(args)
	if err != nil {
		return []string{}, err
	}

	g.parsed = true

	return g.args, nil
}

type argType int

const (
	argEmpty argType = iota
	argNoDash
	argSingleDash
	argDoubleDash
	argSingleDashOpt
	argDoubleDashOpt
)

func classifyArg(arg string) argType {
	switch {
	case arg == "":
		return argEmpty
	case arg == "-":
		return argSingleDash
	case arg == "--":
		return argDoubleDash
	case arg[0] != '-':
		return argNoDash
	case len(arg) > 2 && arg[0:2] == "--":
		return argDoubleDashOpt
	case len(arg) > 1:
		return argSingleDashOpt
	default:
		return argNoDash
	}
}

func (g *Group) shouldStopParsing(arg string, remainingArgs []string) bool {
	switch classifyArg(arg) {
	case argEmpty, argNoDash, argSingleDash:
		// Stop parsing and keep arg in g.args.
		return true
	case argDoubleDash:
		// Stop parsing but drop "--" from g.args.
		g.args = remainingArgs
		return true
	default:
		// Keep parsing and don't change g.args.
		return false
	}
}

func (g *Group) parseByArgType(arg string, args []string) ([]string, error) {
	switch classifyArg(arg) {
	case argSingleDashOpt:
		return g.parseOpt(arg[1:], args)
	case argDoubleDashOpt:
		return g.parseOpt(arg[2:], args)
	default:
		return args, fmt.Errorf("opts: malformed argument: %s", arg)
	}
}

func (g *Group) parse(args []string) error {
	g.args = args

	for len(args) > 0 {
		arg := args[0]
		args = args[1:]

		if g.shouldStopParsing(arg, args) {
			return nil
		}

		var err error
		args, err = g.parseByArgType(arg, args)
		if err != nil {
			return err
		}

		g.args = args
	}

	return nil
}

func (g *Group) parseOpt(arg string, args []string) ([]string, error) {
	name, value, eqFound := strings.Cut(arg, "=")

	opt, ok := g.opts[name]
	if !ok {
		return nil, fmt.Errorf("opts: unknown option --%s", name)
	}

	if eqFound {
		return parseEquals(opt, name, value, arg, args)
	}

	return parseSpaced(opt, name, args)
}

func parseEquals(opt *opt, name, value, arg string, args []string) ([]string, error) {
	if opt.isBool {
		return nil, fmt.Errorf("opts: --%s=%s: boolean options do not accept values", name, value)
	}

	if err := opt.value.set(value); err != nil {
		// Distinguish no value from a bad value.
		if value == "" {
			return nil, fmt.Errorf("opts: --%s: %w", name, ErrMissingValue)
		}

		return nil, &InvalidValueError{
			Option: name,
			Value:  value,
			Err:    err,
		}
	}

	// A string option `--foo=` will not produce an error when calling set.
	// `--foo=` amounts to `--foo=""`, and the empty string is a valid
	// string value. However, for consistency with other option types, we
	// should return an error indicating that there is no value.
	if value == "" && arg[len(arg)-1] == '=' {
		return nil, fmt.Errorf("opts: --%s: %w", name, ErrMissingValue)
	}

	return args, nil
}

func parseSpaced(opt *opt, name string, args []string) ([]string, error) {
	var value string

	switch {
	case opt.isBool:
		value = "true"
	case len(args) > 0:
		value, args = args[0], args[1:]
	default:
		return nil, fmt.Errorf("opts: --%s: %w", name, ErrMissingValue)
	}

	if err := opt.value.set(value); err != nil {
		return nil, &InvalidValueError{
			Option: name,
			Value:  value,
			Err:    err,
		}
	}

	return args, nil
}
