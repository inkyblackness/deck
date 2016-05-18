package release

import (
	"os"
	"path"
)

type directoryReleaseContainer struct {
	path string
}

// NewContainerFromDir returns a release container for given file path.
// An error is returned if the specified directory doesn't exist or any other problem occurred.
func NewContainerFromDir(path string) (container ReleaseContainer, err error) {
	file, err := os.Open(path)

	if file != nil {
		defer file.Close()

		container = &directoryReleaseContainer{path: path}
	}

	return
}

func (container *directoryReleaseContainer) Names() []string {
	dirs := []string{}
	file, _ := os.Open(container.path)

	if file != nil {
		defer file.Close()
		files, _ := file.Readdir(0)

		for _, entry := range files {
			if entry.IsDir() {
				dirs = append(dirs, entry.Name())
			}
		}
	}

	return dirs
}

func (container *directoryReleaseContainer) Get(name string) (rel Release, err error) {
	return ReleaseFromDir(path.Join(container.path, name))
}

func (container *directoryReleaseContainer) New(name string) (rel Release, err error) {
	releasePath := path.Join(container.path, name)
	err = os.Mkdir(releasePath, 0755)

	if err == nil {
		rel, err = ReleaseFromDir(releasePath)
	}

	return
}
