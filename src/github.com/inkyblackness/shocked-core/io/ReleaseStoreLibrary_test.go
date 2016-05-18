package io

import (
	"time"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/chunk/dos"
	"github.com/inkyblackness/shocked-core/release"
)

import (
	check "gopkg.in/check.v1"
)

type ReleaseStoreLibrarySuite struct {
	source  release.Release
	sink    release.Release
	library StoreLibrary
}

var _ = check.Suite(&ReleaseStoreLibrarySuite{})

func (suite *ReleaseStoreLibrarySuite) SetUpTest(c *check.C) {
	suite.source = release.NewMemoryRelease()
	suite.sink = release.NewMemoryRelease()
	suite.library = NewReleaseStoreLibrary(suite.source, suite.sink, 0)
}

func (suite *ReleaseStoreLibrarySuite) createChunkResource(rel release.Release, name string, filler func(consumer chunk.Consumer)) {
	resource, _ := rel.NewResource(name, "")
	writer, _ := resource.AsSink()
	consumer := dos.NewChunkConsumer(writer)
	filler(consumer)
	consumer.Finish()
}

func (suite *ReleaseStoreLibrarySuite) TestChunkStoreIsBackedBySinkIfExisting(c *check.C) {
	suite.createChunkResource(suite.sink, "fromSink.res", func(consumer chunk.Consumer) {
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
	})
	store, err := suite.library.ChunkStore("fromSink.res")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	blockStore := store.Get(res.ResourceID(1))
	c.Check(blockStore.BlockCount(), check.Equals, uint16(1))
}

func (suite *ReleaseStoreLibrarySuite) TestChunkStoreIsBackedBySinkIfExistingInBoth(c *check.C) {
	suite.createChunkResource(suite.source, "stillFromSink.res", func(consumer chunk.Consumer) {})
	suite.createChunkResource(suite.sink, "stillFromSink.res", func(consumer chunk.Consumer) {
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
	})
	store, err := suite.library.ChunkStore("stillFromSink.res")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	blockStore := store.Get(res.ResourceID(1))
	c.Check(blockStore.BlockCount(), check.Equals, uint16(1))
}

func (suite *ReleaseStoreLibrarySuite) TestChunkStoreIsBackedBySourceIfExisting(c *check.C) {
	suite.createChunkResource(suite.source, "fromSource.res", func(consumer chunk.Consumer) {
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
	})
	store, err := suite.library.ChunkStore("fromSource.res")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	blockStore := store.Get(res.ResourceID(1))
	c.Check(blockStore.BlockCount(), check.Equals, uint16(1))
}

func (suite *ReleaseStoreLibrarySuite) TestChunkStoreReturnsEmptyStoreIfNowhereExisting(c *check.C) {
	store, err := suite.library.ChunkStore("empty.res")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	ids := store.IDs()
	c.Check(len(ids), check.Equals, 0)
}

func (suite *ReleaseStoreLibrarySuite) TestModifyingSourceSavesNewSink(c *check.C) {
	suite.createChunkResource(suite.source, "source.res", func(consumer chunk.Consumer) {
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
	})
	store, err := suite.library.ChunkStore("source.res")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)

	store.Del(res.ResourceID(1))

	time.Sleep(100 * time.Millisecond)

	c.Check(suite.sink.HasResource("source.res"), check.Equals, true)
}

func (suite *ReleaseStoreLibrarySuite) TestChunkStoreReturnsSameInstances(c *check.C) {
	suite.createChunkResource(suite.source, "source.res", func(consumer chunk.Consumer) {
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
	})
	store1, err1 := suite.library.ChunkStore("source.res")
	c.Assert(err1, check.IsNil)
	c.Assert(store1, check.NotNil)

	store2, err2 := suite.library.ChunkStore("source.res")
	c.Assert(err2, check.IsNil)

	c.Check(store1, check.Equals, store2)
}
