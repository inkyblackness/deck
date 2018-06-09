package chunks

import "github.com/inkyblackness/res/chunk"

func mapChunk(compressed bool, block []byte) *chunk.Chunk {
	return &chunk.Chunk{
		Compressed:    compressed,
		ContentType:   chunk.Map,
		Fragmented:    false,
		BlockProvider: chunk.MemoryBlockProvider([][]byte{block})}
}
