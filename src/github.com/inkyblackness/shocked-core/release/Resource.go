package release

import "github.com/inkyblackness/res/serial"

// Resource describes one resource file.
type Resource interface {
	// Name returns the unique identifier - the file name of the resource.
	Name() string
	// Path returns the (relative) path for the resource, based on the release's root.
	Path() string
	// AsSource returns an interface for reading the resource.
	AsSource() (serial.SeekingReadCloser, error)
	// AsSink returns an interface for writing the resource.
	AsSink() (serial.SeekingWriteCloser, error)
}
