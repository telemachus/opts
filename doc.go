/*
Package opts creates and parses command line options.

# Basics

An opts [Group] stores a set of options and provides methods to define and
parse them. [NewGroup] creates an empty group. The parameter passed to NewGroup
becomes the [*Group.Name].

Options (aka, flags) are added to a Group by calling type-specific methods on
the Group. These methods come in two forms: long and short. The long versions
take three arguments: a pointer to store the value, a name, and a default
value. The short versions take only the pointer and name arguments since they
default to the zero value for that type. E.g., [*Group.StringZero] defaults to
"" and [*Group.IntZero] defaults to 0. Boolean options are an exception: they
always default to false. Thus, there is only [*Group.Bool] for boolean options.

Valid option names must not be empty, must not begin with "-", and must not
contain whitespace, control characters, quotes, backslashes, or equal signs.
Option definition methods will panic if a name is invalid.

After all options are defined, a call to [*Group.Parse] will attempt to assign
values to options. Parse takes a slice of strings as an argument. This slice
should contain everything on the command line other than the name of the
program. In other words, users should usually pass os.Args[1:], or its
equivalent, to Parse.

Parse will return an error if a user passes an unknown option, if an option is
missing a value, or if a value cannot be parsed as its type. If Parse returns
without error, then the variables associated with options are ready for use. If
Parse returns an error, then those variables are not safe to use.

Parsing stops at the first non-option argument. After successful parsing,
[*Group.Args] returns remaining command-line arguments.

As a usage example, imagine a tool that validates and optionally corrects the
case conventions in one or more files.

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
	if err := og.Parse(os.Args[1:]); err != nil {
		// Handle the error.
	}

	// No error? The values in cfg are safe to use.

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
opt is a string, int, or float, users will pass the obvious thing. E.g.,
"something.toml", "0", or "12.3".

But it may not be obvious what [Group.Date] and [Group.Duration] consider valid
or invalid strings. [Group.Date] opts must be in in RFC 3339 full-date format.
E.g., "2025-12-31" or "2024-02-29". See [cloud.google.com/go/civil.ParseDate]
for details. [Group.Duration] opts must be valid [time.Duration] string. E.g.,
"10ms", "3m2s", or "1h35m9s1ms".  See [time.ParseDuration] for details.
*/
package opts
