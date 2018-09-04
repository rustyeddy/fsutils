package fsutils

import "os"

// Node represents an object in the file system, either a directory
// or a file.
type Node struct {
	Path        string
	os.FileInfo // If this is nil, we are not synced with FS
	*Error
}

// Info returns the file info for this node, that will tell us
// wether this node exists, and if so, all its stats, created,
// size, etc.
func (n *Node) Info() (fi os.FileInfo, fserr *Error) {
	var err error
	if n.FileInfo == nil {
		if n.FileInfo, err = os.Stat(n.Path); err != nil {
			n.Error = NewError(n.Path, err, "os stat info")
		}
	}
	return n.FileInfo, n.Error
}
