package resfile

import (
	"testing"

	"github.com/inkyblackness/res/serial"

	"github.com/inkyblackness/res/chunk"
	"github.com/stretchr/testify/assert"
)

func TestNewWriterReturnsErrorForNilTarget(t *testing.T) {
	writer, err := NewWriter(nil)

	assert.Nil(t, writer, "writer should be nil")
	assert.Equal(t, errTargetNil, err)
}

func TestWriterFinishWithoutAddingChunksCreatesValidFileWithoutChunks(t *testing.T) {
	emptyFileData := emptyResourceFile()
	store := serial.NewByteStore()
	writer, err := NewWriter(store)
	assert.Nil(t, err, "no error expected creating writer")

	err = writer.Finish()
	assert.Nil(t, err, "no error expected finishing writer")
	assert.Equal(t, emptyFileData, store.Data())
}

func TestWriterFinishReturnsErrorWhenAlreadyFinished(t *testing.T) {
	writer, _ := NewWriter(serial.NewByteStore())

	writer.Finish()

	err := writer.Finish()
	assert.Equal(t, errWriterFinished, err)
}

func TestWriterUncompressedSingleBlockChunkCanBeWritten(t *testing.T) {
	data := []byte{0xAB, 0x01, 0xCD, 0x02, 0xEF}
	store := serial.NewByteStore()
	writer, _ := NewWriter(store)
	chunkWriter, err := writer.CreateChunk(chunk.ID(0x1234), chunk.ContentType(0x0A), false)
	assert.Nil(t, err, "no error expected")
	chunkWriter.Write(data)
	writer.Finish()

	result := store.Data()

	var expected []byte
	expected = append(expected, data...)
	expected = append(expected, 0x00, 0x00, 0x00)       // alignment for directory
	expected = append(expected, 0x01, 0x00)             // chunk count
	expected = append(expected, 0x80, 0x00, 0x00, 0x00) // offset to first chunk
	expected = append(expected, 0x34, 0x12)             // chunk ID
	expected = append(expected, 0x05, 0x00, 0x00)       // chunk length (uncompressed)
	expected = append(expected, 0x00)                   // chunk type (uncompressed, single-block)
	expected = append(expected, 0x05, 0x00, 0x00)       // chunk length in file
	expected = append(expected, 0x0A)                   // content type
	assert.Equal(t, expected, result[chunkDirectoryFileOffsetPos+4:])
}

func TestWriterUncompressedFragmentedChunkCanBeWritten(t *testing.T) {
	blockData1 := []byte{0xAB, 0x01, 0xCD}
	blockData2 := []byte{0x11, 0x22, 0x33, 0x44}
	store := serial.NewByteStore()
	writer, _ := NewWriter(store)
	chunkWriter, err := writer.CreateFragmentedChunk(chunk.ID(0x5678), chunk.ContentType(0x0B), false)
	assert.Nil(t, err, "no error expected")
	chunkWriter.CreateBlock().Write(blockData1)
	chunkWriter.CreateBlock().Write(blockData2)
	writer.Finish()

	result := store.Data()

	var expected []byte
	expected = append(expected, 0x02, 0x00)             // number of blocks
	expected = append(expected, 0x0E, 0x00, 0x00, 0x00) // offset to first block
	expected = append(expected, 0x11, 0x00, 0x00, 0x00) // offset to second block
	expected = append(expected, 0x15, 0x00, 0x00, 0x00) // size of chunk
	expected = append(expected, blockData1...)
	expected = append(expected, blockData2...)
	expected = append(expected, 0x00, 0x00, 0x00)       // alignment for directory
	expected = append(expected, 0x01, 0x00)             // chunk count
	expected = append(expected, 0x80, 0x00, 0x00, 0x00) // offset to first chunk
	expected = append(expected, 0x78, 0x56)             // chunk ID
	expected = append(expected, 0x15, 0x00, 0x00)       // chunk length (uncompressed)
	expected = append(expected, 0x02)                   // chunk type
	expected = append(expected, 0x15, 0x00, 0x00)       // chunk length in file
	expected = append(expected, 0x0B)                   // content type
	assert.Equal(t, expected, result[chunkDirectoryFileOffsetPos+4:])
}

func TestWriterUncompressedFragmentedChunkCanBeWrittenWithPaddingForSpecialID(t *testing.T) {
	blockData1 := []byte{0xAB, 0x01, 0xCD}
	blockData2 := []byte{0x11, 0x22, 0x33, 0x44}
	store := serial.NewByteStore()
	writer, _ := NewWriter(store)
	chunkWriter, err := writer.CreateFragmentedChunk(chunk.ID(0x08FD), chunk.ContentType(0x0B), false)
	assert.Nil(t, err, "no error expected")
	chunkWriter.CreateBlock().Write(blockData1)
	chunkWriter.CreateBlock().Write(blockData2)
	writer.Finish()

	result := store.Data()

	var expected []byte
	expected = append(expected, 0x02, 0x00)             // number of blocks
	expected = append(expected, 0x10, 0x00, 0x00, 0x00) // offset to first block
	expected = append(expected, 0x13, 0x00, 0x00, 0x00) // offset to second block
	expected = append(expected, 0x17, 0x00, 0x00, 0x00) // size of chunk
	expected = append(expected, 0x00, 0x00)             // padding
	expected = append(expected, blockData1...)
	expected = append(expected, blockData2...)
	expected = append(expected, 0x00)                   // alignment for directory
	expected = append(expected, 0x01, 0x00)             // chunk count
	expected = append(expected, 0x80, 0x00, 0x00, 0x00) // offset to first chunk
	expected = append(expected, 0xFD, 0x08)             // chunk ID
	expected = append(expected, 0x17, 0x00, 0x00)       // chunk length (uncompressed)
	expected = append(expected, 0x02)                   // chunk type
	expected = append(expected, 0x17, 0x00, 0x00)       // chunk length in file
	expected = append(expected, 0x0B)                   // content type
	assert.Equal(t, expected, result[chunkDirectoryFileOffsetPos+4:])
}

func TestWriterCompressedSingleBlockChunkCanBeWritten(t *testing.T) {
	data := []byte{0x01, 0x02, 0x01, 0x02}
	store := serial.NewByteStore()
	writer, _ := NewWriter(store)
	chunkWriter, err := writer.CreateChunk(chunk.ID(0x1122), chunk.ContentType(0x0C), true)
	assert.Nil(t, err, "no error expected")
	chunkWriter.Write(data)
	writer.Finish()

	result := store.Data()

	var expected []byte
	// 0000 0000|0000 0100|0000 0000|0010 0000|0100 0000|0011 1111|1111 1111
	expected = append(expected, 0x00, 0x04, 0x00, 0x20, 0x40, 0x3F, 0xFF, 0x00) // 14bit words 0x0001 0x0002 0x0100 0x3FFF + trailing 0x00
	expected = append(expected)                                                 // alignment for directory
	expected = append(expected, 0x01, 0x00)                                     // chunk count
	expected = append(expected, 0x80, 0x00, 0x00, 0x00)                         // offset to first chunk
	expected = append(expected, 0x22, 0x11)                                     // chunk ID
	expected = append(expected, 0x04, 0x00, 0x00)                               // chunk length (uncompressed)
	expected = append(expected, 0x01)                                           // chunk type
	expected = append(expected, 0x08, 0x00, 0x00)                               // chunk length in file
	expected = append(expected, 0x0C)                                           // content type
	assert.Equal(t, expected, result[chunkDirectoryFileOffsetPos+4:])
}

func TestWriterCompressedFragmentedChunkCanBeWritten(t *testing.T) {
	blockData1 := []byte{0x01, 0x02, 0x01, 0x02}
	blockData2 := []byte{0x01, 0x02, 0x01, 0x02}
	store := serial.NewByteStore()
	writer, _ := NewWriter(store)
	chunkWriter, err := writer.CreateFragmentedChunk(chunk.ID(0x5544), chunk.ContentType(0x09), true)
	assert.Nil(t, err, "no error expected")
	chunkWriter.CreateBlock().Write(blockData1)
	chunkWriter.CreateBlock().Write(blockData2)
	writer.Finish()

	result := store.Data()

	var expected []byte
	expected = append(expected, 0x02, 0x00)                               // number of blocks
	expected = append(expected, 0x0E, 0x00, 0x00, 0x00)                   // offset to first block
	expected = append(expected, 0x12, 0x00, 0x00, 0x00)                   // offset to second block
	expected = append(expected, 0x16, 0x00, 0x00, 0x00)                   // size of chunk
	expected = append(expected, 0x00, 0x04, 0x00, 0x20, 0x40)             // compressed data, part 1
	expected = append(expected, 0x01, 0x02, 0x00, 0x0B, 0xFF, 0xF0, 0x00) // compressed data, part 2
	expected = append(expected, 0x00, 0x00)                               // alignment for directory
	expected = append(expected, 0x01, 0x00)                               // chunk count
	expected = append(expected, 0x80, 0x00, 0x00, 0x00)                   // offset to first chunk
	expected = append(expected, 0x44, 0x55)                               // chunk ID
	expected = append(expected, 0x16, 0x00, 0x00)                         // chunk length (uncompressed)
	expected = append(expected, 0x03)                                     // chunk type
	expected = append(expected, 0x1A, 0x00, 0x00)                         // chunk length in file
	expected = append(expected, 0x09)                                     // content type
	assert.Equal(t, expected, result[chunkDirectoryFileOffsetPos+4:])
}
