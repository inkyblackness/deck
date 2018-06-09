package chunk

// Provider provides chunks from some storage.
type Provider interface {
	// IDs returns a list of available IDs this provider can provide.
	IDs() []Identifier

	// Provide returns a chunk for the given identifier.
	Chunk(id Identifier) (*Chunk, error)
}
