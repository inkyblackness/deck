package logic

import (
	"github.com/inkyblackness/res/data"

	check "gopkg.in/check.v1"
)

type LevelObjectChainSuite struct {
	chain   *LevelObjectChain
	entries [4]data.LevelObjectPrefix
}

var _ = check.Suite(&LevelObjectChainSuite{})

func (suite *LevelObjectChainSuite) SetUpTest(c *check.C) {
	linkGetter := func(index data.LevelObjectChainIndex) LevelObjectChainLink {
		return &suite.entries[index]
	}
	suite.chain = NewLevelObjectChain(&suite.entries[data.LevelObjectChainStartIndex], linkGetter)
	suite.chain.Initialize(len(suite.entries) - 1)
}

func (suite *LevelObjectChainSuite) TestInitializeResetsAllFields(c *check.C) {
	c.Check(suite.entries[0], check.DeepEquals, data.LevelObjectPrefix{
		LevelObjectTableIndex: 0,
		Next:     0,
		Previous: 1})

	c.Check(suite.entries[1].Previous, check.Equals, uint16(2))
	c.Check(suite.entries[2].Previous, check.Equals, uint16(3))
	c.Check(suite.entries[3].Previous, check.Equals, uint16(0))
}

func (suite *LevelObjectChainSuite) TestAcquireLinkReturnsIndexForNewEntry(c *check.C) {
	index, _ := suite.chain.AcquireLink()

	c.Check(index, check.Not(check.Equals), data.LevelObjectChainIndex(0))
}

func (suite *LevelObjectChainSuite) TestAcquireLinkUpdatesStartWhenEmpty(c *check.C) {
	index, _ := suite.chain.AcquireLink()

	c.Check(suite.entries[0], check.DeepEquals, data.LevelObjectPrefix{
		LevelObjectTableIndex: uint16(index),
		Next:     uint16(index),
		Previous: 2})
}

func (suite *LevelObjectChainSuite) TestAcquireLinkUpdatesEntries(c *check.C) {
	index1, _ := suite.chain.AcquireLink()
	index2, _ := suite.chain.AcquireLink()

	c.Check(suite.entries[0], check.DeepEquals, data.LevelObjectPrefix{
		LevelObjectTableIndex: uint16(index2),
		Next:     uint16(index1),
		Previous: 3})

	c.Check(suite.entries[index1], check.DeepEquals, data.LevelObjectPrefix{
		LevelObjectTableIndex: 0,
		Next:     uint16(index2),
		Previous: uint16(data.LevelObjectChainStartIndex)})

	c.Check(suite.entries[index2], check.DeepEquals, data.LevelObjectPrefix{
		LevelObjectTableIndex: 0,
		Next:     uint16(data.LevelObjectChainStartIndex),
		Previous: uint16(index1)})
}

func (suite *LevelObjectChainSuite) TestAcquireLinkReturnsErrorWhenExhausted(c *check.C) {
	suite.chain.AcquireLink()
	suite.chain.AcquireLink()
	suite.chain.AcquireLink()

	_, err := suite.chain.AcquireLink()

	c.Check(err, check.NotNil)
}

func (suite *LevelObjectChainSuite) TestReleaseLinkPutsEntryBackOnAvailablePool(c *check.C) {
	index, _ := suite.chain.AcquireLink()

	suite.chain.ReleaseLink(index)

	c.Check(suite.entries[index].Previous, check.Not(check.Equals), uint16(data.LevelObjectChainStartIndex))
	c.Check(suite.entries[0].Previous, check.Equals, uint16(index))
}

func (suite *LevelObjectChainSuite) TestReleaseLinkRestoresUpdatesStartPointer(c *check.C) {
	prevIndex, _ := suite.chain.AcquireLink()
	index, _ := suite.chain.AcquireLink()

	suite.chain.ReleaseLink(index)

	c.Check(suite.entries[0].LevelObjectTableIndex, check.Equals, uint16(prevIndex))
}

func (suite *LevelObjectChainSuite) TestReleaseLinkRestoresNeighbours(c *check.C) {
	prevIndex, _ := suite.chain.AcquireLink()
	index, _ := suite.chain.AcquireLink()
	nextIndex, _ := suite.chain.AcquireLink()

	suite.chain.ReleaseLink(index)

	c.Check(suite.entries[prevIndex].Next, check.Equals, uint16(nextIndex))
	c.Check(suite.entries[nextIndex].Previous, check.Equals, uint16(prevIndex))
}
