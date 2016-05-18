package chunk

// BlockStore is a random access block container.
type BlockStore interface {
	BlockHolder

	// SetBlockData sets the data for the requested block index.
	SetBlockData(block uint16, data []byte)
}
