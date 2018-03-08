package resfile

import (
	"io"
)

type blockReaderFunc func(index int) (io.Reader, error)

// ChunkReader provides meta information as well as reader access to its contained blocks.
type ChunkReader struct {
	fragmented  bool
	contentType ContentType
	compressed  bool
	blockCount  int
	blockReader blockReaderFunc
}

// Fragmented describes how many blocks can be expected.
// Unfragmented chunks have exactly one block, fragmented chunks have zero, one, or more blocks.
func (reader *ChunkReader) Fragmented() bool {
	return reader.fragmented
}

// ContentType describes the nature of the data within the chunk - the format of the blocks.
func (reader *ChunkReader) ContentType() ContentType {
	return reader.contentType
}

// Compressed returns true if the data is to be serialized in compressed form
// in the resource file.
func (reader *ChunkReader) Compressed() bool {
	return reader.compressed
}

// BlockCount returns the number of available blocks in this chunk.
// Unfragmented chunks will always have exactly one block.
func (reader *ChunkReader) BlockCount() int {
	return reader.blockCount
}

// Block returns the reader for the identified block.
// Each call returns a new reader instance.
// Data provided by this reader is always uncompressed.
func (reader *ChunkReader) Block(index int) (io.Reader, error) {
	return reader.blockReader(index)
}
