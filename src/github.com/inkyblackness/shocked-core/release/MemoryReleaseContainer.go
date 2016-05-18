package release

import (
	"fmt"
)

type memoryReleaseContainer struct {
	releases map[string]Release
}

// NewMemoryReleaseContainer returns a ReleaseContainer for releases only in memory.
func NewMemoryReleaseContainer() ReleaseContainer {
	return &memoryReleaseContainer{releases: make(map[string]Release)}
}

// Names returns the list of currently known releases.
func (container *memoryReleaseContainer) Names() []string {
	names := make([]string, 0, len(container.releases))
	for name := range container.releases {
		names = append(names, name)
	}

	return names
}

// Get returns the release with given name, or an error if not possible.
func (container *memoryReleaseContainer) Get(name string) (rel Release, err error) {
	rel, existing := container.releases[name]
	if !existing {
		err = fmt.Errorf("Not found")
	}

	return
}

// New creates a new release with given name and returns it, or an error if not possible.
func (container *memoryReleaseContainer) New(name string) (rel Release, err error) {
	rel, existing := container.releases[name]
	if !existing {
		testingRel := NewMemoryRelease()
		rel = testingRel
		container.releases[name] = testingRel
	} else {
		rel = nil
		err = fmt.Errorf("Release exists")
	}

	return
}
