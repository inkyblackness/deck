package release

import (
	"fmt"

	"github.com/inkyblackness/res/serial"
)

type memoryResource struct {
	name string
	path string

	data       []byte
	readLocks  int
	writeLocks int
}

func NewMemoryResource(name string, path string, data []byte) Resource {
	res := &memoryResource{
		name: name,
		path: path,
		data: data}

	return res
}

// Name returns the unique identifier - the file name of the resource.
func (res *memoryResource) Name() string {
	return res.name
}

// Path returns the (relative) path for the resource, based on the release's root.
func (res *memoryResource) Path() string {
	return res.path
}

// AsSource returns an interface for reading the resource.
func (res *memoryResource) AsSource() (buf serial.SeekingReadCloser, err error) {
	if res.writeLocks == 0 {
		res.readLocks++
		buf = serial.NewByteStoreFromData(res.data, func([]byte) { res.readLocks-- })
	} else {
		err = fmt.Errorf("Cannot open for reading")
	}

	return
}

// AsSink returns an interface for writing the resource.
func (res *memoryResource) AsSink() (buf serial.SeekingWriteCloser, err error) {
	if res.readLocks == 0 && res.writeLocks == 0 {
		res.writeLocks++
		buf = serial.NewByteStoreFromData(nil, func(data []byte) {
			res.data = data
			res.writeLocks--
		})
	} else {
		err = fmt.Errorf("Cannot open for writing")
	}

	return
}
