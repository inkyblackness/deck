package store

import (
	"github.com/inkyblackness/res"

	check "gopkg.in/check.v1"
)

type ProviderBackedSuite struct {
}

var _ = check.Suite(&ProviderBackedSuite{})

func (suite *ProviderBackedSuite) TestIDsReturnsIDsFromProvider_WhenStoreUnchanged(c *check.C) {
	testing := NewTestingResource()
	testing.Consume(res.ResourceID(1), emptyBlockHolder())
	testing.Consume(res.ResourceID(2), emptyBlockHolder())

	backed := NewProviderBacked(testing, func() {})
	result := backed.IDs()

	c.Check(result, check.DeepEquals, []res.ResourceID{res.ResourceID(1), res.ResourceID(2)})
}

func (suite *ProviderBackedSuite) TestIDsReturnsIDWithoutDeleted_WhenChunkDeleted(c *check.C) {
	testing := NewTestingResource()
	testing.Consume(res.ResourceID(1), emptyBlockHolder())
	testing.Consume(res.ResourceID(2), emptyBlockHolder())

	backed := NewProviderBacked(testing, func() {})
	backed.Del(res.ResourceID(1))
	result := backed.IDs()

	c.Check(result, check.DeepEquals, []res.ResourceID{res.ResourceID(2)})
}

func (suite *ProviderBackedSuite) TestIDsReturnsIDWithNewID_WhenChunkAdded(c *check.C) {
	testing := NewTestingResource()
	backed := NewProviderBacked(testing, func() {})
	backed.Put(res.ResourceID(4), emptyBlockHolder())
	result := backed.IDs()

	c.Check(result, check.DeepEquals, []res.ResourceID{res.ResourceID(4)})
}

func (suite *ProviderBackedSuite) TestIDsReturnsUnchangedArray_WhenExistingChunkOverwritten(c *check.C) {
	testing := NewTestingResource()
	testing.Consume(res.ResourceID(5), emptyBlockHolder())
	backed := NewProviderBacked(testing, func() {})
	backed.Put(res.ResourceID(5), emptyBlockHolder())
	result := backed.IDs()

	c.Check(result, check.DeepEquals, []res.ResourceID{res.ResourceID(5)})
}

func (suite *ProviderBackedSuite) TestChunksCanBeOverwrittenAndRemoved(c *check.C) {
	testing := NewTestingResource()
	testing.Consume(res.ResourceID(5), emptyBlockHolder())
	backed := NewProviderBacked(testing, func() {})
	backed.Put(res.ResourceID(5), emptyBlockHolder())
	backed.Del(res.ResourceID(5))
	result := backed.IDs()

	c.Check(result, check.DeepEquals, []res.ResourceID{})
}

func (suite *ProviderBackedSuite) TestChunksCanBeOverwrittenRemovedAndAddedAgain(c *check.C) {
	testing := NewTestingResource()
	testing.Consume(res.ResourceID(5), emptyBlockHolder())
	backed := NewProviderBacked(testing, func() {})
	backed.Put(res.ResourceID(5), emptyBlockHolder())
	backed.Del(res.ResourceID(5))
	backed.Put(res.ResourceID(5), emptyBlockHolder())
	result := backed.IDs()

	c.Check(result, check.DeepEquals, []res.ResourceID{res.ResourceID(5)})
}

func (suite *ProviderBackedSuite) TestModifiedCallback_WhenChunkDeleted(c *check.C) {
	testing := NewTestingResource()
	testing.Consume(res.ResourceID(1), emptyBlockHolder())
	testing.Consume(res.ResourceID(2), emptyBlockHolder())

	called := false
	backed := NewProviderBacked(testing, func() { called = true })
	backed.Del(res.ResourceID(2))

	c.Check(called, check.Equals, true)
}

func (suite *ProviderBackedSuite) TestModifiedCallbackNotCalled_WhenUnknownChunkDeleted(c *check.C) {
	testing := NewTestingResource()
	testing.Consume(res.ResourceID(1), emptyBlockHolder())
	testing.Consume(res.ResourceID(2), emptyBlockHolder())

	called := false
	backed := NewProviderBacked(testing, func() { called = true })
	backed.Del(res.ResourceID(3))

	c.Check(called, check.Equals, false)
}

func (suite *ProviderBackedSuite) TestModifiedCallback_WhenChunkAdded(c *check.C) {
	testing := NewTestingResource()
	called := false
	backed := NewProviderBacked(testing, func() { called = true })
	backed.Put(res.ResourceID(7), emptyBlockHolder())

	c.Check(called, check.Equals, true)
}

func (suite *ProviderBackedSuite) TestModifiedCallback_WhenChunkOverwritten(c *check.C) {
	testing := NewTestingResource()
	testing.Consume(res.ResourceID(7), emptyBlockHolder())
	called := false
	backed := NewProviderBacked(testing, func() { called = true })
	backed.Put(res.ResourceID(7), emptyBlockHolder())

	c.Check(called, check.Equals, true)
}

func (suite *ProviderBackedSuite) TestModifiedCallback_WhenChunkBlockDataModified(c *check.C) {
	testing := NewTestingResource()
	testing.Consume(res.ResourceID(7), emptyBlockHolder())
	called := false
	backed := NewProviderBacked(testing, func() { called = true })
	backed.Get(res.ResourceID(7)).SetBlockData(0, []byte{0x01})

	c.Check(called, check.Equals, true)
}

func (suite *ProviderBackedSuite) TestGetReturnsNil_WhenUnknownIDRequested(c *check.C) {
	testing := NewTestingResource()
	backed := NewProviderBacked(testing, func() {})
	result := backed.Get(res.ResourceID(10))

	c.Check(result, check.IsNil)
}

func (suite *ProviderBackedSuite) TestGetReturnsBlockStore_WhenKnownIDRequested(c *check.C) {
	testing := NewTestingResource()
	testing.Consume(res.ResourceID(20), emptyBlockHolder())
	backed := NewProviderBacked(testing, func() {})
	result := backed.Get(res.ResourceID(20))

	c.Check(result, check.Not(check.IsNil))
}

func (suite *ProviderBackedSuite) TestGetReturnsProperBlockStore_WhenKnownIDRequested(c *check.C) {
	testing := NewTestingResource()
	holderA := randomBlockHolder(3)
	holderB := randomBlockHolder(4)
	testing.Consume(res.ResourceID(20), holderA)
	testing.Consume(res.ResourceID(30), holderB)
	backed := NewProviderBacked(testing, func() {})
	result20 := backed.Get(res.ResourceID(20))
	result30 := backed.Get(res.ResourceID(30))

	c.Check(result20.BlockCount(), check.Equals, holderA.BlockCount())
	c.Check(result30.BlockCount(), check.Equals, holderB.BlockCount())
}

func (suite *ProviderBackedSuite) TestGetReturnsSameBlockStoreAsBefore_WhenKnownIDRequested(c *check.C) {
	testing := NewTestingResource()
	holder := randomBlockHolder(1)
	testing.Consume(res.ResourceID(20), holder)
	backed := NewProviderBacked(testing, func() {})
	data := []byte{0x11, 0xAA, 0xBB}
	first := backed.Get(res.ResourceID(20))
	first.SetBlockData(0, data)

	second := backed.Get(res.ResourceID(20))

	c.Check(second.BlockData(0), check.DeepEquals, data)
}

func (suite *ProviderBackedSuite) TestGetReturnsNil_WhenDeletedIDRequested(c *check.C) {
	testing := NewTestingResource()
	testing.Consume(res.ResourceID(20), emptyBlockHolder())
	backed := NewProviderBacked(testing, func() {})

	backed.Del(res.ResourceID(20))

	result := backed.Get(res.ResourceID(20))
	c.Check(result, check.IsNil)
}

func (suite *ProviderBackedSuite) TestGetReturnsBlockStore_WhenNewIDRequested(c *check.C) {
	testing := NewTestingResource()
	backed := NewProviderBacked(testing, func() {})
	backed.Put(res.ResourceID(20), emptyBlockHolder())
	result := backed.Get(res.ResourceID(20))

	c.Check(result, check.Not(check.IsNil))
}
