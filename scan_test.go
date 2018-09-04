package fsutils

import (
	"testing"
)

func ScanTest(t *testing.T) {
	var dirs []string
	dirs = append(dirs, "/")
	Scan(dirs)
}
