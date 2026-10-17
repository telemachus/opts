# opts version history

## v0.5.0

+ `Parse` provides relaxed parsing.
  It returns leftover positional arguments plus any error from parsing.
  Previously, this was the behavior of `ParseKnown`.
+ `Parse` now returns `nil` instead of `[]string{}` on error.
+ `ParseAll` provides strict parsing.
  It returns `*UnexpectedArgumentsError` if there are leftover positional arguments.
  Previously, this was the behavior of `Parse`.
+ `ParseKnown` no longer exists.
+ opts now requires at least Go 1.26.
+ Uint options accept only decimal numbers.
  This matches the existing behavior of Int options.
+ Float64 options accept only decimal notation.
  opts now rejects NaN, Inf, hex floats, and underscores.
  opts still accepts exponents (e.g., 1e3).
+ opts copies leftover arguments returned by Parse.
  They do not share memory with the caller's slice.
  (The same applies to `*UnexpectedArgumentsError.Args`.)
+ `--name=` is now valid for string options.
  The `name` option will be set to `""`.
  (This gives consistency with `--name ""`, which was already valid.)
+ Error messages now report options with one or two dashes as the user typed them.
+ A failed parse no longer leaves options partly set.
  The code now lives up to its original promise that failed parses do not affect options.

For changes before v0.5.0, see the git history or tags v0.1.0 to v0.4.0.
