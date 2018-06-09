package chunk

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func verifyBlockContent(t *testing.T, provider BlockProvider, index int, expected []byte) {
	reader, err := provider.Block(index)
	assert.Nil(t, err, "No error expected for index %d", index)
	assert.NotNil(t, reader, "Reader expected for index %d", index)
	if reader != nil {
		data, dataErr := ioutil.ReadAll(reader)
		assert.Nil(t, dataErr, "No error expected reading data from index %d", index)
		assert.Equal(t, expected, data, "Proper data expected from index %d", index)
	}
}

func verifyBlockError(t *testing.T, provider BlockProvider, index int) {
	_, err := provider.Block(index)
	assert.NotNil(t, err, "Error expected for index %d", index)
}

func TestChunkRefersToBlockProviderByDefault(t *testing.T) {
	chunk := &Chunk{BlockProvider: MemoryBlockProvider([][]byte{{0x01}, {0x02, 0x02}})}

	assert.Equal(t, 2, chunk.BlockCount())
	verifyBlockContent(t, chunk, 0, []byte{0x01})
	verifyBlockContent(t, chunk, 1, []byte{0x02, 0x02})
}

func TestChunkBlockReturnsErrorOnInvalidIndex(t *testing.T) {
	chunk := &Chunk{BlockProvider: MemoryBlockProvider(nil)}

	verifyBlockError(t, chunk, -1)
	verifyBlockError(t, chunk, 0)
	verifyBlockError(t, chunk, 1)
	verifyBlockError(t, chunk, 2)
}

func TestChunkBlockReturnsErrorForDefaultObject(t *testing.T) {
	var chunk Chunk

	assert.Equal(t, 0, chunk.BlockCount())
	verifyBlockError(t, chunk, -1)
	verifyBlockError(t, chunk, 0)
	verifyBlockError(t, chunk, 1)
}

func TestChunkCanBeExtendedWithBlocks(t *testing.T) {
	var chunk Chunk

	chunk.SetBlock(0, []byte{0x10})
	chunk.SetBlock(2, []byte{0x20, 0x20})
	assert.Equal(t, 3, chunk.BlockCount())
	verifyBlockContent(t, chunk, 0, []byte{0x10})
	verifyBlockContent(t, chunk, 1, []byte{})
	verifyBlockContent(t, chunk, 2, []byte{0x20, 0x20})
}

func TestChunkDefaultsToProviderWhenNoExtensionOverridesIt(t *testing.T) {
	chunk := &Chunk{BlockProvider: MemoryBlockProvider([][]byte{{0x01}, {0x02, 0x02}})}

	chunk.SetBlock(0, []byte{0xA0})
	assert.Equal(t, 2, chunk.BlockCount())
	verifyBlockContent(t, chunk, 0, []byte{0xA0})
	verifyBlockContent(t, chunk, 1, []byte{0x02, 0x02})
}

func TestChunkPanicsForNegativeBlockIndex(t *testing.T) {
	var chunk Chunk

	assert.Panics(t, func() { chunk.SetBlock(-1, []byte{}) }, "Panic expected")
}
