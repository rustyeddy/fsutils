package fsutils

// FSError represents a FileSystem Error, which includes
// the Path that caused the error.  Smart directory usage
type Error struct {
	Path    string
	Err     error
	Comment string
}

// FSError
func NewError(path string, err error, comment string) (fserr *Error) {
	fserr = &Error{
		Path:    path,
		Err:     err,
		Comment: comment,
	}
	return fserr
}
