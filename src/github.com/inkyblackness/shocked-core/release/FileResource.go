package release

import (
	"os"
	"path/filepath"

	"github.com/inkyblackness/res/serial"
)

type fileResource struct {
	name         string
	basePath     string
	relativePath string

	filename string
}

func newFileResource(name string, basePath string, relativePath string, filename string) Resource {
	resource := &fileResource{
		name:         name,
		basePath:     basePath,
		relativePath: relativePath,
		filename:     filename}

	return resource
}

func (resource *fileResource) Name() string {
	return resource.name
}

func (resource *fileResource) Path() string {
	return resource.relativePath
}

func (resource *fileResource) AsSource() (serial.SeekingReadCloser, error) {
	return os.Open(filepath.Join(resource.basePath, resource.relativePath, resource.filename))
}

func (resource *fileResource) AsSink() (serial.SeekingWriteCloser, error) {
	os.MkdirAll(filepath.Join(resource.basePath, resource.relativePath), os.FileMode(0755))
	return os.Create(filepath.Join(resource.basePath, resource.relativePath, resource.filename))
}
