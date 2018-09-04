package fsutils

// Keep Some Status along the way
type FileStats struct {
	Basedir     string
	TotalCount  int
	FileCount   int
	DirCount    int
	Unknown     int
	LocalSize   int64 // Just files in this directory
	SubdirsSize int64 // Size of all subdirs combined (not this dir)
}

func (fs *FileStats) TotalSize() int64 {
	return fs.LocalSize + fs.SubdirsSize
}
