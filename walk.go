package fsutils

import (
	"os"
	"path/filepath"
	"sync"
)

type Walker struct {
	Basedir  string
	FileChan chan os.FileInfo
	Stats
}

// NewWalker will create a new directory walker for the given path
func NewWalker(path string) *Walker {
	w := &Walker{
		Basedir:  path,
		FileChan: make(chan os.FileInfo),
	}
	return w
}

// WalkDir does a recursive walk down a directory, sending
// filesizes over the sizeChan channel.
func WalkDir(path string, n *sync.WaitGroup, fiChan chan<- os.FileInfo) {

	// Make sure our wait group is decremented before this
	// function returns
	defer n.Done()

	// Loop each entry and create more subdir searches.  Making
	// sure the waitgroup is updated properly
	for _, entry := range Dirlist(path) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(path, entry.Name())
			go WalkDir(subdir, n, fiChan)
		} else {
			fiChan <- entry
		}
	}
}
