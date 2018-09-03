package main

import (
	"flag"
	"os"

	"github.com/rustyeddy/fsutils"
)

var (
	pattern = flag.String("pattern", "", "Match this pattern regexp or glob")
	glob    = flag.Bool("glob", true, "Treat match as a glob (*.go, ..) ")
	action  = flag.String("action", "scan", "Actions to perform on dir")
)

// Actions are copy, move, remove and xlate

func main() {
	flag.Parse()

	err := fsutils.scan(os.Args[1:])
	fsutils.fatalError(err)
}
