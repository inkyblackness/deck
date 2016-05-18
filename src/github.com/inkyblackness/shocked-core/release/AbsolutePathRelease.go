package release

import (
	"fmt"
	"os"
	"strings"
)

type absolutePathRelease struct {
	paths []string
}

// FromAbsolutePaths returns a Release that provides resources exclusively from
// the provided paths.
func FromAbsolutePaths(paths []string) (release Release, err error) {
	var wrongPaths []string

	for _, path := range paths {
		info, pathErr := os.Stat(path)

		if (pathErr != nil) || !info.IsDir() {
			wrongPaths = append(wrongPaths, path)
		}
	}

	if len(wrongPaths) > 0 {
		err = fmt.Errorf("Paths not valid: %v", wrongPaths)
	} else {
		release = &absolutePathRelease{paths: paths}
	}

	return
}

func (release *absolutePathRelease) resolve(name string) (foundPath string, filename string) {

	for i := 0; i < len(release.paths) && filename == ""; i++ {
		path := release.paths[i]
		dir, _ := os.Open(path)

		if dir != nil {
			defer dir.Close()
			files, err := dir.Readdir(0)

			if err == nil {
				for _, file := range files {
					if !file.IsDir() && strings.ToLower(file.Name()) == name {
						foundPath = path
						filename = file.Name()
					}
				}
			}
		}
	}

	return
}

func (release *absolutePathRelease) HasResource(name string) bool {
	_, filename := release.resolve(name)

	return filename != ""
}

func (release *absolutePathRelease) GetResource(name string) (res Resource, err error) {
	absolutePath, filename := release.resolve(name)

	if filename != "" {
		res = newFileResource(name, "", absolutePath, filename)
	} else {
		err = fmt.Errorf("Resource not found")
	}

	return
}

func (release *absolutePathRelease) NewResource(name string, path string) (res Resource, err error) {
	_, filename := release.resolve(name)

	if filename == "" {
		known := false

		for _, knownPath := range release.paths {
			if knownPath == path {
				known = true
			}
		}
		if known {
			res = newFileResource(name, "", path, name)
		} else {
			err = fmt.Errorf("Unknown path")
		}
	} else {
		err = fmt.Errorf("Resource exists")
	}

	return
}
