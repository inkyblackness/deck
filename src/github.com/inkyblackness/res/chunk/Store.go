package chunk

import "github.com/inkyblackness/res"

// Store represents a dynamically accessible container of chunks.
type Store interface {
	// IDs returns a list of available IDs this store currently contains.
	IDs() []res.ResourceID

	// Get returns a chunk for the given identifier.
	Get(id res.ResourceID) BlockStore

	// Del ensures the chunk with given identifier is deleted.
	Del(id res.ResourceID)

	// Put (re-)assigns an identifier with data. If no chunk with given ID exists,
	// then it is created. Existing chunks are overwritten with the provided data.
	Put(id res.ResourceID, holder BlockHolder)
}
