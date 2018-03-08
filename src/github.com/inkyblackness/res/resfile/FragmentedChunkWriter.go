package resfile

import (
	"io"

	"github.com/inkyblackness/res/resfile/compression"
	"github.com/inkyblackness/res/serial"
)

// FragmentedChunkWriter writes a chunk with zero, one, or more blocks.
// Multiple blocks can be created and then written concurrently. Only when the
// chunk is finished, the blocks are finalized.
type FragmentedChunkWriter struct {
	target *serial.PositioningEncoder

	compressed      bool
	dataPaddingSize int
	blockStores     []*serial.ByteStore
	blockWriter     []*BlockWriter
}

// CreateBlock provides a new, dedicated writer for a new block.
func (writer *FragmentedChunkWriter) CreateBlock() *BlockWriter {
	store := serial.NewByteStore()
	blockWriter := &BlockWriter{target: serial.NewEncoder(store), finisher: func() {}}
	writer.blockStores = append(writer.blockStores, store)
	writer.blockWriter = append(writer.blockWriter, blockWriter)
	return blockWriter
}

func (writer *FragmentedChunkWriter) finish() (length uint32) {
	var unpackedSize uint32
	blockCount := len(writer.blockStores)
	writer.target.Code(uint16(blockCount))
	offset := 2 + (blockCount+1)*4 + writer.dataPaddingSize
	for index, store := range writer.blockStores {
		unpackedSize += writer.blockWriter[index].finish()
		writer.target.Code(uint32(offset))
		offset += len(store.Data())
	}
	writer.target.Code(uint32(offset))
	unpackedSize += writer.target.CurPos() + uint32(writer.dataPaddingSize)

	writer.writeBlocks()

	return unpackedSize
}

func (writer *FragmentedChunkWriter) writeBlocks() {
	var targetWriter io.Writer = writer.target
	targetFinisher := func() {}
	if writer.compressed {
		compressor := compression.NewCompressor(targetWriter)
		targetWriter = compressor
		targetFinisher = func() { compressor.Close() } // nolint: errcheck
	}
	for i := 0; i < writer.dataPaddingSize; i++ {
		targetWriter.Write([]byte{0x00}) // nolint: errcheck
	}
	for _, store := range writer.blockStores {
		targetWriter.Write(store.Data()) // nolint: errcheck
	}
	targetFinisher()
}
