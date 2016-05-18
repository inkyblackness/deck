package release

import (
	"fmt"
)

// StaticReleaseContainer is not capable of creating new releases. It provides
// only the initially registered releases.
type StaticReleaseContainer struct {
	releases map[string]Release
}

// NewStaticReleaseContainer returns a new instance of a ReleaseContainer.
func NewStaticReleaseContainer(releases map[string]Release) ReleaseContainer {
	container := &StaticReleaseContainer{
		releases: make(map[string]Release)}

	for key, release := range releases {
		container.releases[key] = release
	}

	return container
}

// Names implements the ReleaseContainer interface
func (container *StaticReleaseContainer) Names() []string {
	names := make([]string, 0, len(container.releases))

	for key := range container.releases {
		names = append(names, key)
	}

	return names
}

// Get implements the ReleaseContainer interface
func (container *StaticReleaseContainer) Get(name string) (rel Release, err error) {
	rel, existing := container.releases[name]

	if !existing {
		err = fmt.Errorf("Unknown release [%s]", name)
	}

	return
}

// New implements the ReleaseContainer interface
func (container *StaticReleaseContainer) New(name string) (Release, error) {
	return nil, fmt.Errorf("Static container does not support creating new releases")
}
