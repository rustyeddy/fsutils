package fsutils

import (
	"fmt"
	"io/ioutil"
	"os"
)

var (
	MaxConcurrentDirs int // max con
	CountingSemaphore chan struct{}
)

// Create our semaphore based on number of concurrent dirs
func init() {
	// XXX - Must set after flags are parsed in main() after
	// this init function has been run.
	MaxConcurrentDirs = 20
}

// GetSortedEntries takes a path string and returns three arrays
// (lists) of: regular files, directories and other (pipes, perm denied, etc.)
func GetSortedDirlist(path string) (files, dirs, other []os.FileInfo) {
	entries := Dirlist(path)
	if entries == nil {
		return nil, nil, nil
	}
	files, dirs, other = SortDirlist(entries)
	return files, dirs, other
}

// GetEntries converts a path string to an []os.FileInfo each FileInfo
// represents one the the "paths" children.  They can be files,
// [sub]Directories or "other" dependending on the respective file type.
// Directories may be used for deeper search (or not), files may be
// used to information gathering, translation, copy, move or delte, etc.
func Dirlist(path string) (entries []os.FileInfo) {
	var err error
	if CountingSemaphore == nil {
		CountingSemaphore = make(chan struct{}, MaxConcurrentDirs)
	}

	// Limit the number of dir lists we can have running at a time
	CountingSemaphore <- struct{}{}
	defer func() { <-CountingSemaphore }()

	entries, err = ioutil.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "walkDir failed to read %s %v\n", path, err)
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