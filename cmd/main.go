package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

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

	roots := getRootDirs(flag.Args())

	// Create the size channel to report file sizes, simply gather
	// sizes and total them (also count the number of files)
	fiChan := make(chan os.FileInfo)
	var n sync.WaitGroup

	// Create go routines for all roots provided from the
	for _, root := range roots {
		n.Add(1)
		go fsutils.WalkDir(root, &n, fiChan)
	}

	// Wait for all walk functions to complete then
	// close the sizeChan
	go func() {
		n.Wait()
		close(fiChan)
	}()

	// Create a ticker to update the user of progress.  Verbose
	// if true will cause the ticker to emit the scan summary
	// at that point.
	tick := fsutils.CreateTicker(500*time.Millisecond, *verbose)

	// Loop until the sizeChan is closed. The break out of
	// the loop
	var stats fsutils.Stats
loop:
	for {
		select {
		case fi, ok := <-fiChan:
			if !ok {
				break loop // sizeChan was closed ...
			}
			stats.Update(fi)
		case <-tick:
			fmt.Println(stats.String())
		}
	}
	// The final Print usage
	fmt.Printf("Total of %s\n", stats.String())
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
