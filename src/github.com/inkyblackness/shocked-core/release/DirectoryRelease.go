package release

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type directoryRelease struct {
	basePath string
}

// PathDepthLimit specifies how many nested directories are searched for resources.
const PathDepthLimit = 2

// ReleaseFromDir returns a release instance if the given path specifies an existing directory.
func ReleaseFromDir(path string) (release Release, err error) {
	info, err := os.Stat(path)

	if err == nil {
		if info.IsDir() {
			release = &directoryRelease{basePath: path}
		} else {
			err = fmt.Errorf("Not a directory")
		}
	}

	return
}

func (release *directoryRelease) resolve(level int, relPath string, name string) (foundRelPath string, filename string) {
	dir, _ := os.Open(path.Join(release.basePath, relPath))

	if dir != nil {
		defer dir.Close()
		files, err := dir.Readdir(0)

		if err == nil {
			for _, file := range files {
				if file.IsDir() {
					if level < PathDepthLimit && filename == "" {
						foundRelPath, filename = release.resolve(level+1, path.Join(relPath, file.Name()), name)
					}
				} else if strings.ToLower(file.Name()) == name {
					foundRelPath = relPath
					filename = file.Name()
				}
			}
		}
	}

	return
}

func (release *directoryRelease) HasResource(name string) bool {
	_, filename := release.resolve(0, ".", name)

	return len(filename) > 0
}

func (release *directoryRelease) GetResource(name string) (res Resource, err error) {
	relativePath, filename := release.resolve(0, ".", name)

	if filename != "" {
		res = newFileResource(name, release.basePath, relativePath, filename)
	} else {
		err = fmt.Errorf("Resource not found")
	}

	return
}

func (release *directoryRelease) NewResource(name string, path string) (res Resource, err error) {
	_, filename := release.resolve(0, ".", name)

	if filename == "" {
		res = newFileResource(name, release.basePath, path, name)
	} else {
		err = fmt.Errorf("Resource not found")
	}

	return
}
