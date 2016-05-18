package chunk

import "github.com/inkyblackness/res"

// Provider provides chunks
type Provider interface {
	// IDs returns a list of available IDs this provider can provide
	IDs() []res.ResourceID

	// Provide returns a chunk for the given identifier
	Provide(id res.ResourceID) BlockHolder
}
