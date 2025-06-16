/*
Package opts creates and parses command line options.

An opts [*Group] stores a set of options and provides access to their methods.
[opts.NewGroup] creates an empty group. The parameter passed to NewGroup
becomes the [*Group.Name].

As an example, imagine a tool that validates and optionally corrects the case
conventions in one or more files.

	og := opts.NewGroup("caser")

Add options (aka, flags) to the Group by calling type-specific methods on the
Group. These methods come in two forms: long and short. The long versions take
three arguments: a pointer of the appropriate type to hold the value, a name,
and a default value. The short versions take only the pointer and name
arguments since they always default to the zero value for that type. E.g.,
[*Group.StringZero] defaults to "" while [*Group.IntZero] defaults to 0.

Boolean options are an exception: they always default to false. Thus, there is
only [*Group.Bool] for boolean options.

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
	og.UintZero(&cfg.level, "verbosity")
	og.Bool(&cfg.dryRun, "dry-run")
	og.Bool(&cfg.write, "write")

[*Group.Parse] should be called after all options are defined. If Parse returns
without error, then the variables associated with options are ready for use. If
Parse returns an error, then the variables associated with those options are
not safe to use. Parse takes a slice of strings as an argument. This slice
should contain everything on the command line other than the name of the
program. In other words, users should usually pass os.Args[1:] to Parse.
*/
package opts
