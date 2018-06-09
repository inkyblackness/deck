package chunk

type nullProvider struct{}

// NullProvider returns a Provider instance that is empty.
// It contains no IDs and will not provide any chunk.
func NullProvider() Provider {
	return &nullProvider{}
}

func (*nullProvider) IDs() []Identifier {
	return nil
}

func (*nullProvider) Chunk(id Identifier) (*Chunk, error) {
	return nil, ErrChunkDoesNotExist(id)
}
