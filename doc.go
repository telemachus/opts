/*
Package opts creates and parses command line options.

# Basics

An opts [Group] stores a set of options and provides methods to define and
parse them. [NewGroup] creates an empty group. The parameter passed to NewGroup
becomes the [*Group.Name].

Options (aka, flags) are added to a Group by calling type-specific definition
methods. These methods come in two forms: one that requires an explicit
default value and one that uses the type's zero value as the default. E.g.,
[*Group.String] versus [*Group.StringZero]. Boolean options are an exception:
they always default to false, so there is only one method, [*Group.Bool].

Valid option names must not be empty, must not begin with "-", and must not
contain whitespace, control characters, quotes, backslashes, or equal signs.
Option definition methods will panic if a name is invalid.

After all options are defined, Opts provides two methods, [*Group.Parse] and
[*Group.ParseKnown], to assign values to defined options. Both methods take
a slice of strings as an argument, and both methods return a slice of strings
and an error. The slice passed to these methods should contain everything on
the command line other than the name of the program. In other words, users
should usually pass os.Args[1:], or its equivalent, to Parse or ParseKnown. The
slice that the methods return will contain any arguments remaining in the
original slice after parsing has ended.

Both parsing methods will return an error if a user passes an undefined option,
if a non-boolean option is missing a value, if a boolean option includes
a value, or if a value cannot be parsed as its type. If a parsing method
returns without error, then the variables associated with options are ready for
use. If a parsing method returns an error, then those variables are not safe to
use.

Parse is strict, returning [ErrUnexpectedArgs] if any non-option arguments
remain. ParseKnown is relaxed and does not return an error in this situation.
Both methods return the slice of leftover arguments, but only Parse treats
a non-empty slice as an error. The relaxed behavior of ParseKnown is necessary
for programs that accept positional arguments after the options.

As an example, imagine a tool that validates and optionally corrects the case
conventions in one or more files. Since such a tool expects one or more
filenames after options have been set, it makes sense to use ParseKnown.

	og := opts.NewGroup("caser")

	cfg := struct {
		rcfile     string
		convention string
		strictness uint
		verbosity  uint
		dryRun     bool
		write      bool
	}{}

	og.String(&cfg.rcfile, "rcfile", "caser.ini")
	og.String(&cfg.convention, "convention", "camel")
	og.Uint(&cfg.strictness, "strictness", 3)
	og.UintZero(&cfg.verbosity, "verbosity")
	og.Bool(&cfg.dryRun, "dry-run")
	og.Bool(&cfg.write, "write")

	// Later...
	remaining, err := og.ParseKnown(os.Args[1:])
	if err != nil {
		// Handle the error.
	}

	// If there is no error, the values in cfg are ready to use and
	// remaining contains the names of files to check.

# Command Line Syntax

The syntax for options is largely the same as for Go's flag library.

	-option    // one dash is accepted
	--option   // two dashes are accepted
	-option=x  // non-boolean flags only
	-option x  // non-boolean flags only

On the command line, options can begin with one or two dashes; they are
equivalent during parsing. As such, there is no distinction between long and
short options. This means there is no way to stack options. That is, `-abc` is
always read as one option, named "abc", rather than `-a -b -c`. Boolean options
do not accept arguments; they are switches. All boolean options are initially
false. If a boolean option is present on the command line, the option's value
is set to true.

Although the library does not distinguish long from short options when parsing,
it can provide users a short and a long option for use on the command line or
in scripts.

	cfg := struct {
		helpWanted    bool
		versionWanted bool
	}{}

	og.Bool(&cfg.helpWanted, "help")
	og.Bool(&cfg.helpWanted, "h")
	og.Bool(&cfg.versionWanted, "version")
	og.Bool(&cfg.versionWanted, "V")

# Valid Command Line Strings

For most types, it will be clear what a valid string will look like. If an
opt is a string, int, uint, or float, users will pass the obvious thing. E.g.,
"something.toml", "0", or "12.3".

But it may not be obvious what [Group.Date] and [Group.Duration] consider valid
or invalid strings. [Group.Date] opts must be in in RFC 3339 full-date format:
YYYY-MM-DD. E.g., "2025-12-31" or "2024-02-29". [Group.Duration] options must
be valid [time.Duration] string. E.g., "10ms", "3m2s", or "1h35m9s1ms". For
details, see [time.ParseDuration].
*/
package opts
