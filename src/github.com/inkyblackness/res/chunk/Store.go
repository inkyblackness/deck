package chunk

// Store represents a dynamically accessible container of chunks.
type Store interface {
	Provider

	// Del removes the chunk with given identifier from the store.
	// Trying to remove an unknown identifier has no effect.
	Del(id Identifier)

	// Put (re-)assigns an identifier with data. If no chunk with given ID exists,
	// then it is created. Existing chunks are overwritten with the provided data.
	Put(id Identifier, chunk *Chunk)
}
