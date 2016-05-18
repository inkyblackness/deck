package release

// ReleaseContainer can contain any number of releases and provides access to them.
type ReleaseContainer interface {
	// Names returns the list of currently known releases.
	Names() []string
	// Get returns the release with given name, or an error if not possible.
	Get(name string) (Release, error)
	// New creates a new release with given name and returns it, or an error if not possible.
	New(name string) (Release, error)
}
