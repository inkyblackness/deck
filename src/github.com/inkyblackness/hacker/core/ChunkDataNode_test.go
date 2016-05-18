package core

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"

	check "gopkg.in/check.v1"
)

type ChunkDataNodeSuite struct {
	parentNode DataNode

	chunkDataNode DataNode
}

var _ = check.Suite(&ChunkDataNodeSuite{})

func (suite *ChunkDataNodeSuite) TestInfoReturnsListOfAvailableBlockCountAndContentType(c *check.C) {
	holder := chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}, []byte{}})
	suite.chunkDataNode = newChunkDataNode(suite.parentNode, res.ResourceID(0x0200), holder)

	result := suite.chunkDataNode.Info()

	c.Check(result, check.Equals, "Content type: 0x00\nAvailable blocks: 2\nChunk TypeID: 0x00 (Basic)\n")
}

func (suite *ChunkDataNodeSuite) TestResolveReturnsDataNodeForKnownID(c *check.C) {
	holder := chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}, []byte{}})
	suite.chunkDataNode = newChunkDataNode(suite.parentNode, res.ResourceID(0x0200), holder)

	result := suite.chunkDataNode.Resolve("1")

	c.Assert(result, check.NotNil)
	c.Check(result.ID(), check.Equals, "1")
}
