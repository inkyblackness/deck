package io

import (
	"io/ioutil"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"

	"gopkg.in/check.v1"
)

type DynamicChunkStoreSuite struct {
}

var _ = check.Suite(&DynamicChunkStoreSuite{})

func (suite *DynamicChunkStoreSuite) SetUpTest(c *check.C) {
}

func (suite *DynamicChunkStoreSuite) createChunkProvider(filler func(chunk.Store)) chunk.Provider {
	store := chunk.NewProviderBackedStore(chunk.NullProvider())
	filler(store)

	return store
}

func (suite *DynamicChunkStoreSuite) aChunk() *chunk.Chunk {
	return &chunk.Chunk{
		BlockProvider: chunk.MemoryBlockProvider([][]byte{{}})}
}

func (suite *DynamicChunkStoreSuite) TestIDsReturnsResultFromWrapped(c *check.C) {
	provider := suite.createChunkProvider(func(consumer chunk.Store) {
		consumer.Put(chunk.ID(1), suite.aChunk())
	})

	wrappedStore := chunk.NewProviderBackedStore(provider)
	store := NewDynamicChunkStore(wrappedStore, func() {})

	ids := store.IDs()
	c.Check(len(ids), check.Equals, 1)
}

func (suite *DynamicChunkStoreSuite) TestDelDeletesFromWrapped(c *check.C) {
	provider := suite.createChunkProvider(func(consumer chunk.Store) {
		consumer.Put(chunk.ID(1), suite.aChunk())
	})

	wrappedStore := chunk.NewProviderBackedStore(provider)
	store := NewDynamicChunkStore(wrappedStore, func() {})

	store.Del(res.ResourceID(1))

	ids := store.IDs()
	c.Check(len(ids), check.Equals, 0)
}

func (suite *DynamicChunkStoreSuite) TestPutInsertsToWrapped(c *check.C) {
	provider := suite.createChunkProvider(func(consumer chunk.Store) {})

	wrappedStore := chunk.NewProviderBackedStore(provider)
	store := NewDynamicChunkStore(wrappedStore, func() {})

	store.Put(res.ResourceID(1), suite.aChunk())

	ids := store.IDs()
	c.Check(len(ids), check.Equals, 1)
}

func (suite *DynamicChunkStoreSuite) TestGetReturnsBlockFromWrapped(c *check.C) {
	provider := suite.createChunkProvider(func(consumer chunk.Store) {
		consumer.Put(chunk.ID(1), suite.aChunk())
	})

	wrappedStore := chunk.NewProviderBackedStore(provider)
	store := NewDynamicChunkStore(wrappedStore, func() {})

	holder := store.Get(res.ResourceID(1))

	c.Check(holder, check.NotNil)
}

func (suite *DynamicChunkStoreSuite) TestGetReturnsNilIfWrappedDoesntHaveIt(c *check.C) {
	provider := suite.createChunkProvider(func(consumer chunk.Store) {})

	wrappedStore := chunk.NewProviderBackedStore(provider)
	store := NewDynamicChunkStore(wrappedStore, func() {})

	holder := store.Get(res.ResourceID(2))

	c.Check(holder, check.IsNil)
}

func (suite *DynamicChunkStoreSuite) TestBlockHolderModifiesWrapped(c *check.C) {
	provider := suite.createChunkProvider(func(consumer chunk.Store) {
		consumer.Put(chunk.ID(1), suite.aChunk())
	})

	wrappedStore := chunk.NewProviderBackedStore(provider)
	store := NewDynamicChunkStore(wrappedStore, func() {})

	data := []byte{0x01, 0x02}
	holder := store.Get(res.ResourceID(1))
	holder.SetBlockData(0, data)

	wrappedChunk, _ := wrappedStore.Chunk(chunk.ID(1))
	wrappedReader, _ := wrappedChunk.Block(0)
	wrappedData, _ := ioutil.ReadAll(wrappedReader)
	c.Check(wrappedData, check.DeepEquals, data)
}

func (suite *DynamicChunkStoreSuite) TestSwapReplacesWrapped(c *check.C) {
	provider1 := suite.createChunkProvider(func(consumer chunk.Store) {
		consumer.Put(chunk.ID(1), suite.aChunk())
	})
	provider2 := suite.createChunkProvider(func(consumer chunk.Store) {
		consumer.Put(res.ResourceID(2), suite.aChunk())
	})

	testStore := NewDynamicChunkStore(chunk.NewProviderBackedStore(provider1), func() {})
	testStore.Swap(func(oldStore chunk.Store) chunk.Store {
		return chunk.NewProviderBackedStore(provider2)
	})

	ids := testStore.IDs()
	values := make([]uint16, len(ids))
	for index, id := range ids {
		values[index] = id.Value()
	}
	c.Check(values, check.DeepEquals, []uint16{2})
}

func (suite *DynamicChunkStoreSuite) TestBlockHolderModifiesWrappedAfterSwap(c *check.C) {
	provider1 := suite.createChunkProvider(func(consumer chunk.Store) {
		consumer.Put(res.ResourceID(1), suite.aChunk())
	})
	provider2 := suite.createChunkProvider(func(consumer chunk.Store) {
		consumer.Put(res.ResourceID(1), suite.aChunk())
	})

	testStore := NewDynamicChunkStore(chunk.NewProviderBackedStore(provider1), func() {})
	holder := testStore.Get(res.ResourceID(1))

	wrapped2 := chunk.NewProviderBackedStore(provider2)
	testStore.Swap(func(oldStore chunk.Store) chunk.Store { return wrapped2 })

	data := []byte{0x01, 0x02}
	holder.SetBlockData(0, data)

	wrappedChunk, _ := wrapped2.Chunk(chunk.ID(1))
	wrappedReader, _ := wrappedChunk.Block(0)
	wrappedData, _ := ioutil.ReadAll(wrappedReader)
	c.Check(wrappedData, check.DeepEquals, data)
}
