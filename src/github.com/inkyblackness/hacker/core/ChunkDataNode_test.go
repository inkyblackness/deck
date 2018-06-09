package core

import (
	"github.com/inkyblackness/res/chunk"

	"gopkg.in/check.v1"
)

type ChunkDataNodeSuite struct {
	parentNode DataNode

	chunkDataNode DataNode
}

var _ = check.Suite(&ChunkDataNodeSuite{})

func (suite *ChunkDataNodeSuite) TestInfoReturnsListOfAvailableBlockCountAndContentType(c *check.C) {
	holder := &chunk.Chunk{
		ContentType:   chunk.Palette,
		BlockProvider: chunk.MemoryBlockProvider([][]byte{[]byte{}, []byte{}})}
	suite.chunkDataNode = newChunkDataNode(suite.parentNode, chunk.ID(0x0200), holder)

	result := suite.chunkDataNode.Info()

	c.Check(result, check.Equals, ""+
		"Content type: 0x00\n"+
		"Compressed: false\n"+
		"Fragmented: false\n"+
		"Available blocks: 2\n")
}

func (suite *ChunkDataNodeSuite) TestResolveReturnsDataNodeForKnownID(c *check.C) {
	holder := &chunk.Chunk{
		ContentType:   chunk.Palette,
		BlockProvider: chunk.MemoryBlockProvider([][]byte{[]byte{}, []byte{}})}
	suite.chunkDataNode = newChunkDataNode(suite.parentNode, chunk.ID(0x0200), holder)

	result := suite.chunkDataNode.Resolve("1")

	c.Assert(result, check.NotNil)
	c.Check(result.ID(), check.Equals, "1")
}
