package chunk

import (
	"bytes"
	"fmt"
	"io"
)

// Chunk provides meta information as well as access to its contained blocks.
type Chunk struct {
	// Fragmented tells whether the chunk should be serialized with a directory.
	// Fragmented chunks can have zero, one, or more blocks.
	// Unfragmented chunks always have exactly one block.
	Fragmented bool

	// ContentType describes how the block data shall be interpreted.
	ContentType ContentType

	// Compressed tells whether the data shall be serialized in compressed form.
	Compressed bool

	// BlockProvider is the keeper of original block data.
	// This provider will be referred to if no other data was explicitly set.
	BlockProvider BlockProvider

	blockLimit int
	blocks     map[int][]byte
}

// BlockCount returns the number of available blocks in the chunk.
// Unfragmented chunks will always have exactly one block.
func (chunk Chunk) BlockCount() (count int) {
	count = chunk.providerBlockCount()
	if count < chunk.blockLimit {
		count = chunk.blockLimit
	}
	return
}

func (chunk Chunk) providerBlockCount() (count int) {
	if chunk.BlockProvider != nil {
		count = chunk.BlockProvider.BlockCount()
	}
	return
}

// Block returns the reader for the identified block.
// Each call returns a new reader instance.
// Data provided by this reader is always uncompressed.
func (chunk Chunk) Block(index int) (io.Reader, error) {
	if chunk.blocks != nil {
		data, set := chunk.blocks[index]
		if set {
			return bytes.NewReader(data), nil
		} else if chunk.providerBlockCount() <= index {
			return bytes.NewReader(nil), nil
		}
	}
	if chunk.BlockProvider == nil {
		return nil, fmt.Errorf("no blocks available")
	}
	return chunk.BlockProvider.Block(index)
}

// SetBlock registers new data for a block.
// For any block set this way, the block provider of this chunk will no longer be queried.
func (chunk *Chunk) SetBlock(index int, data []byte) {
	if index < 0 {
		panic(fmt.Errorf("index must be a non-negative value"))
	}
	chunk.ensureBlockMap()
	chunk.blocks[index] = data
	if chunk.blockLimit <= index {
		chunk.blockLimit = index + 1
	}
}

func (chunk *Chunk) ensureBlockMap() {
	if chunk.blocks == nil {
		chunk.blocks = make(map[int][]byte)
	}
}
