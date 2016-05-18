package release

import (
	"fmt"
)

type memoryRelease struct {
	resources map[string]Resource
}

// NewMemoryRelease returns a Release instance that keeps all resources and their data in memory.
func NewMemoryRelease() Release {
	return &memoryRelease{resources: make(map[string]Resource)}
}

// HasResource returns true for a unique resource name if the release contains this resource.
func (rel *memoryRelease) HasResource(name string) bool {
	_, existing := rel.resources[name]

	return existing
}

// GetResource returns the resource identified by given name if existing, or an error otherwise.
func (rel *memoryRelease) GetResource(name string) (res Resource, err error) {
	res, existing := rel.resources[name]
	if !existing {
		err = fmt.Errorf("Not found")
	}

	return
}

// NewResource creates a new resource under given path and returns the instance, or an error on failure.
func (rel *memoryRelease) NewResource(name string, path string) (res Resource, err error) {
	res, existing := rel.resources[name]
	if !existing {
		testingRes := NewMemoryResource(name, path, nil)
		res = testingRes
		rel.resources[name] = testingRes
	} else {
		res = nil
		err = fmt.Errorf("Resource exists")
	}

	return
}
