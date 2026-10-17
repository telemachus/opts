package opts

import (
	"time"
)

// Duration defines a time.Duration option with the specified name and default
// value. The argument d points to a time.Duration variable that will store the
// value of the option. Duration will panic if name is not valid or repeats an
// existing option. On the command line, users must pass a string in Go's
// time.Duration format (e.g., "500ms" or "4h15m18s"). See [time.ParseDuration]
// for details.
func (g *Group) Duration(d *time.Duration, name string, defValue time.Duration) {
	define(g, "Duration", d, name, defValue, time.ParseDuration)
}

// DurationZero is like Duration but with a default value of 0.
func (g *Group) DurationZero(d *time.Duration, name string) {
	g.Duration(d, name, 0)
}
