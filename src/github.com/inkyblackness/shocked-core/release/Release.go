package release

// Release wraps a set of unique resources that together make up one release.
type Release interface {
	// HasResource returns true for a unique resource name if the release contains this resource.
	HasResource(name string) bool
	// GetResource returns the resource identified by given name if existing, or an error otherwise.
	GetResource(name string) (Resource, error)
	// NewResource creates a new resource under given path and returns the instance, or an error on failure.
	NewResource(name string, path string) (Resource, error)
}
