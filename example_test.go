package opts

import (
	"fmt"
	"os"
)

func Example_typical() {
	// This example demonstrates how to create an options group and define
	// options for it.
	cfg := struct {
		rcfile     string
		convention string
		strictness uint
		dryRun     bool
		write      bool
	}{}

	og := NewGroup("caser")
	og.String(&cfg.rcfile, "rcfile", "caser.ini")
	og.String(&cfg.convention, "convention", "camel")
	og.Uint(&cfg.strictness, "strictness", 3)
	og.Bool(&cfg.dryRun, "dry-run")
	og.Bool(&cfg.write, "write")

	args := []string{"--strictness", "5", "--dry-run", "--", "-awful-filename.txt", "file2.go"}
	if err := og.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "problem parsing args (%v): %v\n", args, err)
	}
	fmt.Printf("Config: %+v\n", cfg)
	fmt.Printf("Remaining args: %v\n", og.Args())
	// Output:
	// Config: {rcfile:caser.ini convention:camel strictness:5 dryRun:true write:false}
	// Remaining args: [-awful-filename.txt file2.go]
}
