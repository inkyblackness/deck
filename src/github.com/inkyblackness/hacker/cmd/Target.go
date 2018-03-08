package cmd

// A Target is an evaluation target, capabile of processing commands
type Target interface {
	// Load requests to load data files from two paths.
	Load(path1, path2 string) string
	// Save re-encodes all loaded data and overwrites the corresponding files.
	Save() string
	// Info returns the status of the current node.
	Info() string
	// ChangeDirectory switches the currently active node
	ChangeDirectory(path string) string
	// Dump returns a data dump of the current node
	Dump() string
	// Diff returns the difference of the current node to the source.
	Diff(source string) string
	// Put sets bytes at the given offset
	Put(offset uint32, data []byte) string
}
