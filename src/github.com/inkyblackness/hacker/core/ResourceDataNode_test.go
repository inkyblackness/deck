package core

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"

	check "gopkg.in/check.v1"
)

type ResourceDataNodeSuite struct {
	parentNode  DataNode
	chunkHolder *TestingChunkProvider
	node        DataNode
}

var _ = check.Suite(&ResourceDataNodeSuite{})

func (suite *ResourceDataNodeSuite) SetUpTest(c *check.C) {
	suite.chunkHolder = NewTestingChunkProvider()
}

func (suite *ResourceDataNodeSuite) TestInfoReturnsListOfAvailableChunkIDs(c *check.C) {
	suite.chunkHolder.Consume(res.ResourceID(0x0100), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{}))
	suite.chunkHolder.Consume(res.ResourceID(0x0050), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{}))
	suite.node = NewResourceDataNode(suite.parentNode, "testFile.res", suite.chunkHolder, nil)

	result := suite.node.Info()

	c.Check(result, check.Equals, "ResourceFile: testFile.res\nIDs: 0100 0050")
}

func (suite *ResourceDataNodeSuite) TestResolveReturnsDataNodeForKnownID(c *check.C) {
	suite.chunkHolder.Consume(res.ResourceID(0x0100), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{}))
	suite.chunkHolder.Consume(res.ResourceID(0x0050), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{}))
	suite.node = NewResourceDataNode(suite.parentNode, "testFile.res", suite.chunkHolder, nil)

	result := suite.node.Resolve("0050")

	c.Assert(result, check.NotNil)
	c.Check(result.ID(), check.Equals, "0050")
}

func (suite *ResourceDataNodeSuite) TestIDReturnsFileNameInLowerCase(c *check.C) {
	suite.node = NewResourceDataNode(suite.parentNode, "TESTFILE.RES", suite.chunkHolder, nil)

	c.Check(suite.node.ID(), check.Equals, "testfile.res")
}
