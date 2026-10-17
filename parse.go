package opts

import (
	"fmt"
	"slices"
	"strings"
)

// Parse scans args and sets option values defined in the [Group].
//
// Parse should be called after all options are defined and before any option
// values are used. If Parse returns without error, the Group is considered
// parsed and subsequent calls to Parse or [Group.ParseAll] will return
// [ErrAlreadyParsed].
//
// Parsing stops at the first non-option argument. Parse returns any remaining
// arguments as a new slice, which does not share memory with args. Both '-'
// and '--' are considered non-option arguments, and both stop further
// parsing. However, the returned slice will not contain '--', but it will
// contain '-'. (By convention, many programs treat '-' as stdin, but that is
// up to the calling program to decide and handle.)
//
// If Parse encounters an unknown option, an option without a value, or a value
// that cannot be parsed as its type, it returns an error and the Group remains
// unparsed. No option values are changed, and the caller may retry with
// different arguments.
//
// The slice passed to Parse should not include the program name. If using
// `os.Args` directly, the caller should pass `os.Args[1:]`.
//
// Use Parse if your program accepts positional arguments after the options.
func (g *Group) Parse(args []string) ([]string, error) {
	if g.parsed {
		return nil, fmt.Errorf("opts: option group %q: %w", g.name, ErrAlreadyParsed)
	}

	pa, err := g.prepare(args)
	if err != nil {
		return nil, err
	}

	pa.assign()
	g.parsed = true

	return slices.Clone(pa.rest), nil
}

// ParseAll is like [Group.Parse], except that it requires every argument to be
// an option. If any positional arguments remain after parsing stops, ParseAll
// returns [*UnexpectedArgumentsError]. In that case, as with any other error,
// no option values are changed and the Group remains unparsed.
//
// Use ParseAll if your program does not accept positional arguments.
func (g *Group) ParseAll(args []string) error {
	if g.parsed {
		return fmt.Errorf("opts: option group %q: %w", g.name, ErrAlreadyParsed)
	}

	pa, err := g.prepare(args)
	if err != nil {
		return err
	}

	if len(pa.rest) > 0 {
		return &UnexpectedArgumentsError{Args: slices.Clone(pa.rest)}
	}

	pa.assign()
	g.parsed = true

	return nil
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

// preparedArgs holds the result of preparing args: the arguments left when
// parsing stops and, in order, one function per option that assigns that
// option's value.
type preparedArgs struct {
	rest        []string
	assignments []func()
}

func (pa preparedArgs) assign() {
	for _, assign := range pa.assignments {
		assign()
	}
}

// prepare scans args without assigning any values.
func (g *Group) prepare(args []string) (preparedArgs, error) {
	var assignments []func()

	for len(args) > 0 {
		arg := args[0]

		var optPrefix string
		switch classifyArg(arg) {
		case argEmpty, argNoDash, argSingleDash:
			// Stop parsing and keep arg.
			return preparedArgs{rest: args, assignments: assignments}, nil
		case argDoubleDash:
			// Stop parsing but drop "--".
			return preparedArgs{rest: args[1:], assignments: assignments}, nil
		case argSingleDashOpt:
			optPrefix = arg[:1]
		case argDoubleDashOpt:
			optPrefix = arg[:2]
		}

		optArg := arg[len(optPrefix):]
		if strings.HasPrefix(optArg, "=") {
			// With no name before "=", report the whole argument.
			return preparedArgs{}, fmt.Errorf("opts: %w %q", ErrUnknownOption, arg)
		}

		rest, assign, err := g.parseOpt(optPrefix, optArg, args[1:])
		if err != nil {
			return preparedArgs{}, err
		}

		assignments = append(assignments, assign)
		args = rest
	}

	return preparedArgs{rest: args, assignments: assignments}, nil
}

// parseOpt parses one option. optPrefix is the single or double dash before
// the option, and arg is the rest of it. Errors show the option as the user
// typed it. If the value is the next argument, parseOpt consumes it from args.
func (g *Group) parseOpt(optPrefix, arg string, args []string) ([]string, func(), error) {
	name, value, eqFound := strings.Cut(arg, "=")

	opt, ok := g.opts[name]
	if !ok {
		return nil, nil, fmt.Errorf("opts: %w %q", ErrUnknownOption, optPrefix+name)
	}

	switch {
	case opt.isBool:
		if eqFound {
			return nil, nil, fmt.Errorf("opts: %w %q", ErrBooleanWithValue, optPrefix+name)
		}

		value = "true"
	case eqFound:
		// The value came after "=".
	case len(args) > 0:
		value, args = args[0], args[1:]
	default:
		return nil, nil, fmt.Errorf("opts: %w for %q", ErrMissingValue, optPrefix+name)
	}

	assign, err := opt.value.prepare(value)
	if err != nil {
		// Distinguish no value from a bad value. A string option accepts
		// `--foo=` as the empty string, as it accepts `--foo ""`.
		if eqFound && value == "" {
			return nil, nil, fmt.Errorf("opts: %w for %q", ErrMissingValue, optPrefix+name+"=")
		}

		return nil, nil, &InvalidValueError{
			Option:    name,
			Value:     value,
			Err:       err,
			optPrefix: optPrefix,
		}
	}

	return args, assign, nil
}
