package fsutils

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Walker struct {
	Roots []string
	Stats
	Verbose bool

	sync.WaitGroup
	FiChan  chan os.FileInfo
	DirChan chan string
	Tick    <-chan time.Time
	*log.Logger
}

// NewWalker will create a new directory walker for the given path
func NewWalker(roots []string) *Walker {
	w := &Walker{
		Roots:   roots,
		FiChan:  make(chan os.FileInfo),
		DirChan: make(chan string),
		Logger:  log.New(),
	}
	return w
}

// WalkDir does a recursive walk down a directory, sending
// filesizes over the sizeChan channel.
func (w *Walker) WalkDir(path string) {

	w.Debugln(" walking dir ", path)

	// Make sure our wait group is decremented before this
	// function returns
	defer func() {
		w.Done()
	}()

	// Loop each entry and create more subdir searches.  Making
	// sure the waitgroup is updated properly
	for _, entry := range Dirlist(path) {
		if entry.IsDir() {
			w.Add(1)
			subdir := filepath.Join(path, entry.Name())
			go w.WalkDir(subdir)
			//w.DirChan <- subdir
		} else {
			w.FiChan <- entry
		}
	}
}

func (w *Walker) StartWalking() {
	// Create the size channel to report file sizes, simply gather
	// sizes and total them (also count the number of files)
	// fiChan := make(chan os.FileInfo) // in walker

	// Create go routines for all roots provided from the
	for _, root := range w.Roots {
		w.Add(1)
		/*
			go func() {
				w.DirChan <- root
			}()
		*/
		go w.WalkDir(root)
	}

	// Wait for all walk functions to complete then close the
	// sizeChan, when everything completes we will.
	go func() {
		w.Wait()
		close(w.FiChan)
	}()

	// Create a ticker to update the user of progress.  Verbose
	// if true will cause the ticker to emit the scan summary
	// at that point.
	w.Tick = CreateTicker(500*time.Millisecond, w.Verbose)
	w.StatsLoop()
}

func (w *Walker) StatsLoop() {
	for {
		select {

		case fi, ok := <-w.FiChan:
			if !ok {
				return
			}
			w.Update(fi)

		case path, ok := <-w.DirChan:
			if !ok {
				return
			}
			w.Add(1)
			w.Stats.Dirs++
			go w.WalkDir(path)

		case <-w.Tick:
			// the tick channel will be readable every 1sec (or ...)
			// it prints an update on os.Stdio for the user.  If
			// verbose is false, the tick channel is never written to.
			fmt.Println(w.Stats.String())
		}
	}

	// The final Print usage
	fmt.Printf("Total of %s\n", w.Stats.String())
	//for path, err := range readErrors {
	//	fmt.Printf("\t%-40s: %v\n", path, err)
	//}
}
