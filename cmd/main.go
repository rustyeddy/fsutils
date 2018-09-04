package main

import (
	"flag"
	"fmt"
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
)

func main() {
	flag.Parse()

	// Set the root file-systems for this this search
	roots := flag.Args()
	if len(roots) == 0 {
		// We could default to this directory. or Fail
		// fmt.Fprintf(os.Stderr, "Need arguments to proceed ... ")

		// We will default to the local directory
		roots = []string{"."}
	}

	// Create the size channel to report file sizes, simply gather
	// sizes and total them (also count the number of files)
	sizeChan := make(chan int64)
	var n sync.WaitGroup

	// Create go routines for all roots provided from the
	for _, root := range roots {
		n.Add(1)
		go fsutils.WalkDir(root, &n, sizeChan)
	}

	// Wait for all walk functions to complete then
	// close the sizeChan
	go func() {
		n.Wait()
		close(sizeChan)
	}()

	// Create the tick chan, the channel will effectively
	// be ignored if verbosity is off.
	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}

	// Loop until the sizeChan is closed. The break out of
	// the loop
	var stats Stats
loop:
	for {
		select {
		case size, ok := <-sizeChan:
			if !ok {
				break loop // sizeChan was closed ...
			}
			stats.Files++
			stats.TotalSize += size
		case <-tick:
			printUsage(stats)
		}
	}
	// The final Print usage
	fmt.Printf("Total of ")
	printUsage(stats)
}

func printUsage(s Stats) {
	fmt.Printf("%d files at %.1f GB\n ", s.Files, float64(s.TotalSize)/1e9)
}
