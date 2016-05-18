package chunk

import "github.com/inkyblackness/res"

// BlockHolder represents a list of blocks of binary data.
type BlockHolder interface {
	// Type returns the type of the chunk.
	ChunkType() TypeID

	// ContentType returns the type of the data.
	ContentType() res.DataTypeID

	// BlockCount returns the number of blocks available in the chunk.
	// Flat chunks must contain exactly one block.
	BlockCount() uint16

	// BlockData returns the data for the requested block index.
	BlockData(block uint16) []byte
}
