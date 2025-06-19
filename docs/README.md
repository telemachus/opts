[![Go Reference](https://pkg.go.dev/badge/github.com/telemachus/opts.svg)](https://pkg.go.dev/github.com/telemachus/opts)

# opts: a small and opinionated argument parser for Go

## Basics

An opts `Group` stores a set of options and provides methods to define and parse
them. `NewGroup` creates an empty group. The parameter passed to NewGroup
becomes the `*Group.Name`.

Options (aka, flags) are added to a Group by calling type-specific methods on
the Group. These methods come in two forms: long and short. The long versions
take three arguments: a pointer to store the value, a name, and a default value.
The short versions take only the pointer and name arguments since they default
to the zero value for that type. E.g., `*Group.StringZero` defaults to "" and
`*Group.IntZero` defaults to 0. Boolean options are an exception: they always
default to false. Thus, there is only `*Group.Bool` for boolean options.

Valid option names must not be empty, must not begin with "-", and must not
contain whitespace, control characters, quotes, backslashes, or equal signs.
Option definition methods will panic if a name is invalid.

After all options are defined, a call to `*Group.Parse` will attempt to assign
values to options. Parse takes a slice of strings as an argument. This slice
should contain everything on the command line other than the name of the
program. In other words, users should usually pass os.Args[1:], or its
equivalent, to Parse.

Parse will return an error if a user passes an unknown option, if an option is
missing a value, or if a value cannot be parsed as its type. If Parse returns
without error, then the variables associated with options are ready for use. If
Parse returns an error, then those variables are not safe to use.

Parsing stops at the first non-option argument. After successful parsing,
`*Group.Args` returns remaining command-line arguments.

As a usage example, imagine a tool that validates and optionally corrects the
case conventions in one or more files.

```go
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
```

## Command Line Syntax

The syntax for options is largely the same as for Go's flag library.

```go
-option    // one dash is accepted
--option   // two dashes are accepted
-option=x  // non-boolean flags only
-option x  // non-boolean flags only
```

On the command line, options can begin with one or two dashes; they are
equivalent during parsing. As such, there is no distinction between long and
short options. This means there is no way to stack options. That is, `-abc` is
always read as one option, named "abc", rather than `-a -b -c`. Boolean options
do not accept arguments; they are switches. All boolean options are initially
false. If a boolean option is present on the command line, the option's value is
set to true.

Although the library does not distinguish long from short options when parsing,
it can provide users a short and a long option for use on the command line or in
scripts.

```go
cfg := struct {
	helpWanted    bool
	versionWanted bool
}{}

og.Bool(&cfg.helpWanted, "help")
og.Bool(&cfg.helpWanted, "h")
og.Bool(&cfg.versionWanted, "version")
og.Bool(&cfg.versionWanted, "V")
```

## Opinionated?

Here are some of the key ways that this library is opinionated.

+ Single dash and double dash are identical. When parsing arguments, the library
  treats, e.g., `-help` and `--help` as if they were identical. (In this way,
  the library follows `flag` in Go's standard library.)
+ No (traditional) short options and no automatic binding of long and short
  options.  The library does not distinguish between short options (preceded by
  a single dash, always one letter and, stackable) versus long options (preceded
  by two dashes, more than one letter, not stackable). Although users can bind
  two options to one variable, the library does not provide methods that take
  two options at once and bind them to the same variable. (Again, this is like
  Go's `flag` library.)
+ No automatic usage. Although most option parsing libraries provide ways to
  generate formatted help messages, `opts` does not. It's more work to write
  help messages by hand, but I think the results can be worth it.
+ Booleans. Booleans always default to false, and they never accept arguments.
  They function only as switches: if a boolean option appears on the command
  line, it's value becomes true.
+ Types are limited and not extendable. The library provides options for the
  following types: boolean, duration, float64, int, string, and uint. Users
  cannot extend the types.
