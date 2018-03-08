package resfile

import (
	"bytes"
	"encoding/binary"

	"github.com/inkyblackness/res/serial"
)

func emptyResourceFile() []byte {
	buf := bytes.NewBufferString(headerString)
	headerTrailer := make([]byte, chunkDirectoryFileOffsetPos-buf.Len())
	headerTrailer[0] = commentTerminator

	binary.Write(buf, binary.LittleEndian, headerTrailer)
	dictionaryOffset := uint32(buf.Len() + 4)
	binary.Write(buf, binary.LittleEndian, &dictionaryOffset)

	numberOfChunks := uint16(0)
	firstChunkOffset := uint32(buf.Len())

	binary.Write(buf, binary.LittleEndian, &numberOfChunks)
	binary.Write(buf, binary.LittleEndian, &firstChunkOffset)

	return buf.Bytes()
}

const (
	exampleChunkIDSingleBlockChunk           = ChunkID(0x4000)
	exampleChunkIDSingleBlockChunkCompressed = ChunkID(0x1000)
	exampleChunkIDFragmentedChunk            = ChunkID(0x2000)
	exampleChunkIDFragmentedChunkCompressed  = ChunkID(0x5000)
)

func exampleResourceFile() []byte {
	store := serial.NewByteStore()
	writer, _ := NewWriter(store)

	chunk1, _ := writer.CreateChunk(exampleChunkIDSingleBlockChunk, ContentType(0x01), false)
	chunk1.Write([]byte{0x01, 0x01, 0x01})
	chunk2, _ := writer.CreateChunk(exampleChunkIDSingleBlockChunkCompressed, ContentType(0x02), true)
	chunk2.Write([]byte{0x02, 0x02})
	chunk3, _ := writer.CreateFragmentedChunk(exampleChunkIDFragmentedChunk, ContentType(0x03), false)
	chunk3.CreateBlock().Write([]byte{0x30, 0x30, 0x30, 0x30})
	chunk3.CreateBlock().Write([]byte{0x31, 0x31, 0x31})
	chunk4, _ := writer.CreateFragmentedChunk(exampleChunkIDFragmentedChunkCompressed, ContentType(0x04), true)
	chunk4.CreateBlock().Write([]byte{0x40, 0x40})
	chunk4.CreateBlock().Write([]byte{0x41, 0x41, 0x41, 0x41})
	chunk4.CreateBlock().Write([]byte{0x42})
	writer.Finish()

	return store.Data()
}
