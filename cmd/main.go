package main

import (
	"flag"
	"fmt"

	"github.com/rustyeddy/fsutils"
)

type Stats struct {
	Files     int64
	Dirs      int64
	Others    int64
	TotalSize int64
}

var (
	action  = flag.String("action", "scan", "Actions to perform on dir")
	glob    = flag.Bool("glob", true, "Treat match as a glob (*.go, ..) ")
	pattern = flag.String("pattern", "", "Match this pattern regexp or glob")
	verbose = flag.Bool("verbose", false, "Print progress and other stuff")

	format = flag.String("format", "text", "Output format text, JSON ... ")
)

func main() {
	flag.Parse()

	walker := fsutils.NewWalker(getRootDirs(flag.Args()))
	walker.Verbose = true
	walker.StartWalking()
	fmt.Printf(" Walker %+v \n", walker)
}

// getRootDirs will default to current directory
func getRootDirs(d []string) (roots []string) {
	roots = flag.Args()
	if len(roots) == 0 {
		// We could default to this directory. or Fail
		// fmt.Fprintf(os.Stderr, "Need arguments to proceed ... ")
		// We will default to the local directory
		roots = []string{"."}
	}
	return roots
}
