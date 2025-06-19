package opts

import "errors"

// ErrAlreadyParsed signals an attempt to parse a [Group] that has already been
// successfully parsed.
var ErrAlreadyParsed = errors.New("opts: option group already parsed")
