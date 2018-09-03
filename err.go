package fsutils

// FSError represents a FileSystem Error, which includes
// the Path that caused the error.  Smart directory usage
type FSError struct {
	Path    string
	Err     error
	Comment string
}

// FSError
func NewFSError(path string, err error, comment string) (fserr *FSError) {
	fserr = &FSError{
		Path:    path,
		Err:     err,
		Comment: comment,
	}
	return fserr
}
