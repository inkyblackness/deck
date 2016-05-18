package dos

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/compress/base"
	"github.com/inkyblackness/res/serial"
)

type blockReader struct {
	coder   serial.PositioningCoder
	address *chunkAddress

	blocks [][]byte
}

// Type returns the type of the chunk.
func (reader *blockReader) ChunkType() chunk.TypeID {
	return chunk.TypeID(reader.address.chunkType)
}

// ContentType returns the type of the data.
func (reader *blockReader) ContentType() res.DataTypeID {
	return res.DataTypeID(reader.address.contentType)
}

// BlockCount returns the number of blocks available in the chunk.
// Flat chunks must contain exactly one block.
func (reader *blockReader) BlockCount() uint16 {
	reader.ensureBlocksBuffered()

	return uint16(len(reader.blocks))
}

// BlockData returns the data for the requested block index.
func (reader *blockReader) BlockData(block uint16) []byte {
	reader.ensureBlocksBuffered()

	return reader.blocks[block]
}

func (reader *blockReader) ensureBlocksBuffered() {
	if reader.blocks == nil {
		blockCoder := serial.Coder(reader.coder)

		reader.coder.SetCurPos(reader.address.startOffset)
		if reader.ChunkType().HasDirectory() {
			blockCount := uint16(0)
			firstStartOffset := uint32(0)

			reader.coder.CodeUint16(&blockCount)
			reader.blocks = make([][]byte, blockCount)

			reader.coder.CodeUint32(&firstStartOffset)
			lastStartOffset := firstStartOffset
			for i := uint16(0); i < blockCount; i++ {
				nextStartOffset := uint32(0)
				reader.coder.CodeUint32(&nextStartOffset)
				reader.blocks[i] = make([]byte, nextStartOffset-lastStartOffset)
				lastStartOffset = nextStartOffset
			}
			reader.coder.SetCurPos(reader.address.startOffset + firstStartOffset) // is this true for compressed as well?
		} else {
			reader.blocks = [][]byte{make([]byte, reader.address.uncompressedLength)}
		}

		if reader.ChunkType().IsCompressed() {
			decompressor := base.NewDecompressor(reader.coder)
			blockCoder = serial.NewDecoder(decompressor)
		}

		for _, data := range reader.blocks {
			blockCoder.CodeBytes(data)
		}
	}
}
