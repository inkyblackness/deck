package resfile

import (
	"bytes"
	"testing"

	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/serial"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	target := serial.NewByteStore()
	provider := chunk.NewProviderBackedStore(chunk.NullProvider())
	aChunk := func(compressed bool, contentType chunk.ContentType, fragmented bool, blocks [][]byte) *chunk.Chunk {
		return &chunk.Chunk{
			Compressed:    compressed,
			ContentType:   contentType,
			Fragmented:    fragmented,
			BlockProvider: chunk.MemoryBlockProvider(blocks)}
	}

	provider.Put(chunk.ID(1), aChunk(false, chunk.Bitmap, false, [][]byte{{0x11}}))
	provider.Put(chunk.ID(3), aChunk(false, chunk.Font, true, [][]byte{{0x21}, {0x22, 0x23}}))
	provider.Put(chunk.ID(2), aChunk(true, chunk.Geometry, false, [][]byte{{0x31}}))
	provider.Put(chunk.ID(4), aChunk(true, chunk.Map, true, [][]byte{{0x41}, {0x42, 0x43}}))

	errWrite := Write(target, provider)
	if errWrite != nil {
		assert.Nil(t, errWrite, "no error expected writing")
	}

	reader, errReader := ReaderFrom(bytes.NewReader(target.Data()))
	if errReader != nil {
		assert.Nil(t, errReader, "no error expected reading")
	}

	assert.Equal(t, []chunk.Identifier{chunk.ID(1), chunk.ID(3), chunk.ID(2), chunk.ID(4)}, reader.IDs())
}
