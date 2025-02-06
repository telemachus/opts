# TODO: opts

## Parsing

+ Bools should not accept arguments.
+ Arguments with values should always use `=`.  Thus, they do not interact with
  further args at all.   (Is this actually a good idea?)

## Types

+ Remove Duration?

## Testing

+ Test that boolean options do not take arguments.
+ Add mixed parsing tests.
+ Separate out tests for single and double dash parsing into a single file.
  Then mix some single and some double in all other files, but don't duplicate
  so much.
