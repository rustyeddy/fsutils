package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Stats struct {
	Files     int64
	Dirs      int64
	Others    int64
	TotalSize int64
}

var (
	pattern = flag.String("pattern", "", "Match this pattern regexp or glob")
	glob    = flag.Bool("glob", true, "Treat match as a glob (*.go, ..) ")
	action  = flag.String("action", "scan", "Actions to perform on dir")
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

	// Create the size channel to report sizes
	sizeChan := make(chan int64)

	// Create go routines to walk each of the root filesystems
	// That have been provided on the command line.  NOTE: we are
	// calling the loop in a go routine, allowing us to drop past
	// the loop. However, the go routine still incures the delay
	// of each write in the loop.
	go func() {
		for _, root := range roots {
			walkDir(root, sizeChan)
		}
		close(sizeChan)
	}()

	// Print the results
	var stats Stats
	for size := range sizeChan {
		stats.Files++
		stats.TotalSize += size
	}
	printUsage(stats)
}

func printUsage(s Stats) {
	fmt.Printf("%d files - %.1f GB\n ", s.Files, float64(s.TotalSize)/1e9)
}

// walkDir does a recursive walk down a directory, sending
// filesizes over the sizeChan channel.
func walkDir(path string, sizeChan chan<- int64) {
	for _, entry := range Dirlist(path) {
		if entry.IsDir() {
			subdir := filepath.Join(path, entry.Name())
			walkDir(subdir, sizeChan)
		} else {
			sizeChan <- entry.Size()
		}
	}
}

// GetSortedEntries takes a path string and returns three arrays
// (lists) of: regular files, directories and other (pipes, perm denied, etc.)
func GetSortedDirlist(path string) (files, dirs, other []os.FileInfo) {
	entries := Dirlist(path)
	if entries == nil {
		return nil, nil, nil
	}
	f, d, o := SortDirlist(entries)
	return f, d, o
}

// GetEntries converts a path string to an []os.FileInfo each FileInfo
// represents one the the "paths" children.  They can be files,
// [sub]Directories or "other" dependending on the respective file type.
// Directories may be used for deeper search (or not), files may be
// used to information gathering, translation, copy, move or delte, etc.
func Dirlist(path string) (entries []os.FileInfo) {
	var err error

	entries, err = ioutil.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "waskDir failed to read %s %v\n", path, err)
		return nil
	}
	return entries
}

// SortDirlist walks through the entire directory list, identifies
// each entry as a: file, dir or other grouping them in the
// appropriate list.
func SortDirlist(dirlist []os.FileInfo) (files, dirs, other []os.FileInfo) {
	for _, fi := range dirlist {
		mode := fi.Mode()
		fmt.Printf("\n mode %+v \n", mode)
		switch {
		case mode.IsDir():
			fmt.Printf("  dir: %s\n", fi.Name())
			dirs = append(dirs, fi) // Could send to a new channel?
		case mode.IsRegular():
			fmt.Printf("  reg: %s\n", fi.Name())
			files = append(files, fi)
		default:
			// TODO ~ Categorize the "other" category with a map?
			other = append(other, fi)
			fmt.Printf(" Other: %s perms or ? %+v\n", fi.Name(), mode)
		}
	}
	return files, dirs, other
}
