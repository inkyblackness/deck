package core

import (
	"github.com/inkyblackness/res/chunk"

	"gopkg.in/check.v1"
)

type ResourceDataNodeSuite struct {
	parentNode DataNode
	store      chunk.Store
	node       DataNode
}

var _ = check.Suite(&ResourceDataNodeSuite{})

func (suite *ResourceDataNodeSuite) SetUpTest(c *check.C) {
	suite.store = chunk.NewProviderBackedStore(chunk.NullProvider())
}

func (suite *ResourceDataNodeSuite) aChunk() *chunk.Chunk {
	return &chunk.Chunk{
		ContentType:   chunk.Palette,
		BlockProvider: chunk.MemoryBlockProvider([][]byte{})}
}

func (suite *ResourceDataNodeSuite) TestInfoReturnsListOfAvailableChunkIDs(c *check.C) {
	suite.store.Put(chunk.ID(0x0100), suite.aChunk())
	suite.store.Put(chunk.ID(0x0050), suite.aChunk())
	suite.node = NewResourceDataNode(suite.parentNode, "testFile.res", suite.store, nil)

	result := suite.node.Info()

	c.Check(result, check.Equals, "ResourceFile: testFile.res\nIDs: 0100 0050")
}

func (suite *ResourceDataNodeSuite) TestResolveReturnsDataNodeForKnownID(c *check.C) {
	suite.store.Put(chunk.ID(0x0100), suite.aChunk())
	suite.store.Put(chunk.ID(0x0050), suite.aChunk())
	suite.node = NewResourceDataNode(suite.parentNode, "testFile.res", suite.store, nil)

	result := suite.node.Resolve("0050")

	c.Assert(result, check.NotNil)
	c.Check(result.ID(), check.Equals, "0050")
}

func (suite *ResourceDataNodeSuite) TestIDReturnsFileNameInLowerCase(c *check.C) {
	suite.node = NewResourceDataNode(suite.parentNode, "TESTFILE.RES", suite.store, nil)

	c.Check(suite.node.ID(), check.Equals, "testfile.res")
}
