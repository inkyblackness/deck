package chunks

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
)

// AddStaticChunk adds a static block with given data to the consumer
func AddStaticChunk(consumer chunk.Consumer, chunkID res.ResourceID, block []byte) {
	AddTypedStaticChunk(consumer, chunkID, chunk.BasicChunkType, block)
}

// AddTypedStaticChunk adds a static block with given data and type to the consumer
func AddTypedStaticChunk(consumer chunk.Consumer, chunkID res.ResourceID, chunkType chunk.TypeID, block []byte) {
	blocks := [][]byte{block}
	consumer.Consume(chunkID, chunk.NewBlockHolder(chunkType, res.Map, blocks))
}
