package chunks

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
)

// AddStaticChunk adds a static block with given data to the consumer
func AddStaticChunk(consumer chunk.Consumer, chunkID res.ResourceID, block []byte) {
	blocks := [][]byte{block}
	consumer.Consume(chunkID, chunk.NewBlockHolder(chunk.BasicChunkType, res.Map, blocks))
}
