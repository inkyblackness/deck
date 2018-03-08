package resfile

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReaderFromReturnsErrorForNilSource(t *testing.T) {
	reader, err := ReaderFrom(nil)

	assert.Nil(t, reader, "reader should be nil")
	assert.Equal(t, errSourceNil, err)
}

func TestReaderFromReturnsInstanceOnEmptySource(t *testing.T) {
	source := bytes.NewReader(emptyResourceFile())
	reader, err := ReaderFrom(source)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, reader)
}

func TestReaderFromReturnsErrorOnInvalidHeaderString(t *testing.T) {
	sourceData := emptyResourceFile()
	sourceData[10] = byte("A"[0])

	_, err := ReaderFrom(bytes.NewReader(sourceData))

	assert.Equal(t, errFormatMismatch, err)
}

func TestReaderFromReturnsErrorOnMissingCommentTerminator(t *testing.T) {
	sourceData := emptyResourceFile()
	sourceData[len(headerString)] = byte(0)

	_, err := ReaderFrom(bytes.NewReader(sourceData))

	assert.Equal(t, errFormatMismatch, err)
}

func TestReaderFromReturnsErrorOnInvalidDirectoryStart(t *testing.T) {
	sourceData := emptyResourceFile()
	sourceData[chunkDirectoryFileOffsetPos] = byte(0xFF)
	sourceData[chunkDirectoryFileOffsetPos+1] = byte(0xFF)

	_, err := ReaderFrom(bytes.NewReader(sourceData))

	assert.NotNil(t, err, "error expected")
}

func TestReaderFromCanDecodeExampleResourceFile(t *testing.T) {
	_, err := ReaderFrom(bytes.NewReader(exampleResourceFile()))
	assert.Nil(t, err, "no error expected")
}

func TestReaderIDsReturnsTheStoredChunkIDsInOrderFromFile(t *testing.T) {
	reader, _ := ReaderFrom(bytes.NewReader(exampleResourceFile()))

	assert.Equal(t, []ChunkID{exampleChunkIDSingleBlockChunk, exampleChunkIDSingleBlockChunkCompressed,
		exampleChunkIDFragmentedChunk, exampleChunkIDFragmentedChunkCompressed}, reader.IDs())
}

func TestReaderChunkReturnsErrorForUnknownID(t *testing.T) {
	reader, _ := ReaderFrom(bytes.NewReader(emptyResourceFile()))
	chunkReader, err := reader.Chunk(ChunkID(0x1111))
	assert.Nil(t, chunkReader, "no reader expected")
	assert.NotNil(t, err)
}

func TestReaderChunkReturnsAChunkReaderForKnownID(t *testing.T) {
	reader, _ := ReaderFrom(bytes.NewReader(exampleResourceFile()))
	chunkReader, err := reader.Chunk(exampleChunkIDSingleBlockChunk)
	assert.Nil(t, err, "no error expected")
	assert.NotNil(t, chunkReader)
}

func TestReaderChunkReturnsChunkWithMetaInformation(t *testing.T) {
	reader, _ := ReaderFrom(bytes.NewReader(exampleResourceFile()))
	info := func(chunkID ChunkID, name string, expected interface{}) string {
		return fmt.Sprintf("Chunk 0x%04X should have %v = %v", chunkID.Value(), name, expected)
	}
	verifyChunk := func(chunkID ChunkID, fragmented bool, contentType ContentType, compressed bool) {
		chunkReader, _ := reader.Chunk(chunkID)
		assert.Equal(t, fragmented, chunkReader.Fragmented(), info(chunkID, "fragmented", fragmented))
		assert.Equal(t, contentType, chunkReader.ContentType(), info(chunkID, "contentType", contentType))
		assert.Equal(t, compressed, chunkReader.Compressed(), info(chunkID, "compressed", compressed))
	}
	verifyChunk(exampleChunkIDSingleBlockChunk, false, ContentType(0x01), false)
	verifyChunk(exampleChunkIDSingleBlockChunkCompressed, false, ContentType(0x02), true)
	verifyChunk(exampleChunkIDFragmentedChunk, true, ContentType(0x03), false)
	verifyChunk(exampleChunkIDFragmentedChunkCompressed, true, ContentType(0x04), true)
}

func TestReaderChunkWithUncompressedSingleBlockContent(t *testing.T) {
	reader, _ := ReaderFrom(bytes.NewReader(exampleResourceFile()))
	chunkReader, _ := reader.Chunk(exampleChunkIDSingleBlockChunk)

	assert.Equal(t, 1, chunkReader.BlockCount())
	verifyBlockContent(t, chunkReader, 0, []byte{0x01, 0x01, 0x01})
}

func TestReaderChunkWithCompressedSingleBlockContent(t *testing.T) {
	reader, _ := ReaderFrom(bytes.NewReader(exampleResourceFile()))
	chunkReader, _ := reader.Chunk(exampleChunkIDSingleBlockChunkCompressed)

	assert.Equal(t, 1, chunkReader.BlockCount())
	verifyBlockContent(t, chunkReader, 0, []byte{0x02, 0x02})
}

func TestReaderChunkWithUncompressedFragmentedContent(t *testing.T) {
	reader, _ := ReaderFrom(bytes.NewReader(exampleResourceFile()))
	chunkReader, _ := reader.Chunk(exampleChunkIDFragmentedChunk)

	assert.Equal(t, 2, chunkReader.BlockCount())
	verifyBlockContent(t, chunkReader, 0, []byte{0x30, 0x30, 0x30, 0x30})
	verifyBlockContent(t, chunkReader, 1, []byte{0x31, 0x31, 0x31})
}

func TestReaderChunkWithCompressedFragmentedContent(t *testing.T) {
	reader, _ := ReaderFrom(bytes.NewReader(exampleResourceFile()))
	chunkReader, _ := reader.Chunk(exampleChunkIDFragmentedChunkCompressed)

	assert.Equal(t, 3, chunkReader.BlockCount())
	verifyBlockContent(t, chunkReader, 0, []byte{0x40, 0x40})
	verifyBlockContent(t, chunkReader, 1, []byte{0x41, 0x41, 0x41, 0x41})
	verifyBlockContent(t, chunkReader, 2, []byte{0x42})
}

func verifyBlockContent(t *testing.T, chunkReader *ChunkReader, blockIndex int, expected []byte) {
	blockReader, readerErr := chunkReader.Block(blockIndex)
	assert.Nil(t, readerErr, "error retrieving reader")
	require.NotNil(t, blockReader, "reader is nil")
	data, dataErr := ioutil.ReadAll(blockReader)
	assert.Nil(t, dataErr, "no error expected reading data")
	assert.Equal(t, expected, data)
}
