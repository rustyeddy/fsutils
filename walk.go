package fsutils

import (
	"fmt"
	"path/filepath"
	"sync"
)

func p(s string) { fmt.Println(s) }

// WalkDir does a recursive walk down a directory, sending
// filesizes over the sizeChan channel.
func WalkDir(path string, n *sync.WaitGroup, sizeChan chan<- int64) {

	// Make sure our wait group is decremented before this
	// function returns
	defer n.Done()

	// Loop each entry and create more subdir searches.  Making
	// sure the waitgroup is updated properly
	for _, entry := range Dirlist(path) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(path, entry.Name())
			go WalkDir(subdir, n, sizeChan)
		} else {
			sizeChan <- entry.Size()
		}
	}
}

// WalkDir does a recursive walk of a filesystem sending
// FileInfo back on the given channel.
/*
func WalkDirChan(path string, n *sync.WaitGroup, fileChan chan<- os.FileInof) {
}
*/
