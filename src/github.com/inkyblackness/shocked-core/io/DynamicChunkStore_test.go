package io

import (
	"bytes"
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/chunk/dos"
	"github.com/inkyblackness/res/chunk/store"
	"github.com/inkyblackness/res/serial"
)

import (
	check "gopkg.in/check.v1"
)

type DynamicChunkStoreSuite struct {
}

var _ = check.Suite(&DynamicChunkStoreSuite{})

func (suite *DynamicChunkStoreSuite) SetUpTest(c *check.C) {
}

func (suite *DynamicChunkStoreSuite) createChunkProvider(filler func(consumer chunk.Consumer)) chunk.Provider {
	store := serial.NewByteStore()
	consumer := dos.NewChunkConsumer(store)
	filler(consumer)
	consumer.Finish()

	provider, _ := dos.NewChunkProvider(bytes.NewReader(store.Data()))

	return provider
}

func (suite *DynamicChunkStoreSuite) TestIDsReturnsResultFromWrapped(c *check.C) {
	provider := suite.createChunkProvider(func(consumer chunk.Consumer) {
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
	})

	wrappedStore := store.NewProviderBacked(provider, func() {})
	store := NewDynamicChunkStore(wrappedStore)

	ids := store.IDs()
	c.Check(len(ids), check.Equals, 1)
}

func (suite *DynamicChunkStoreSuite) TestDelDeletesFromWrapped(c *check.C) {
	provider := suite.createChunkProvider(func(consumer chunk.Consumer) {
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
	})

	wrappedStore := store.NewProviderBacked(provider, func() {})
	store := NewDynamicChunkStore(wrappedStore)

	store.Del(res.ResourceID(1))

	ids := store.IDs()
	c.Check(len(ids), check.Equals, 0)
}

func (suite *DynamicChunkStoreSuite) TestPutInsertsToWrapped(c *check.C) {
	provider := suite.createChunkProvider(func(consumer chunk.Consumer) {})

	wrappedStore := store.NewProviderBacked(provider, func() {})
	store := NewDynamicChunkStore(wrappedStore)

	store.Put(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))

	ids := store.IDs()
	c.Check(len(ids), check.Equals, 1)
}

func (suite *DynamicChunkStoreSuite) TestGetReturnsBlockFromWrapped(c *check.C) {
	provider := suite.createChunkProvider(func(consumer chunk.Consumer) {
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
	})

	wrappedStore := store.NewProviderBacked(provider, func() {})
	store := NewDynamicChunkStore(wrappedStore)

	holder := store.Get(res.ResourceID(1))

	c.Check(holder, check.NotNil)
}

func (suite *DynamicChunkStoreSuite) TestGetReturnsNilIfWrappedDoesntHaveIt(c *check.C) {
	provider := suite.createChunkProvider(func(consumer chunk.Consumer) {})

	wrappedStore := store.NewProviderBacked(provider, func() {})
	store := NewDynamicChunkStore(wrappedStore)

	holder := store.Get(res.ResourceID(2))

	c.Check(holder, check.IsNil)
}

func (suite *DynamicChunkStoreSuite) TestBlockHolderModifiesWrapped(c *check.C) {
	provider := suite.createChunkProvider(func(consumer chunk.Consumer) {
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
	})

	wrappedStore := store.NewProviderBacked(provider, func() {})
	store := NewDynamicChunkStore(wrappedStore)

	data := []byte{0x01, 0x02}
	holder := store.Get(res.ResourceID(1))
	holder.SetBlockData(0, data)

	c.Check(wrappedStore.Get(res.ResourceID(1)).BlockData(0), check.DeepEquals, data)
}

func (suite *DynamicChunkStoreSuite) TestSwapReplacesWrapped(c *check.C) {
	provider1 := suite.createChunkProvider(func(consumer chunk.Consumer) {
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
	})
	provider2 := suite.createChunkProvider(func(consumer chunk.Consumer) {
		consumer.Consume(res.ResourceID(2), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
	})

	testStore := NewDynamicChunkStore(store.NewProviderBacked(provider1, func() {}))
	testStore.Swap(func(oldStore chunk.Store) chunk.Store {
		return store.NewProviderBacked(provider2, func() {})
	})

	c.Check(testStore.IDs(), check.DeepEquals, []res.ResourceID{res.ResourceID(2)})
}

func (suite *DynamicChunkStoreSuite) TestBlockHolderModifiesWrappedAfterSwap(c *check.C) {
	provider1 := suite.createChunkProvider(func(consumer chunk.Consumer) {
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
	})
	provider2 := suite.createChunkProvider(func(consumer chunk.Consumer) {
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
	})

	testStore := NewDynamicChunkStore(store.NewProviderBacked(provider1, func() {}))
	holder := testStore.Get(res.ResourceID(1))

	wrapped2 := store.NewProviderBacked(provider2, func() {})
	testStore.Swap(func(oldStore chunk.Store) chunk.Store { return wrapped2 })

	data := []byte{0x01, 0x02}
	holder.SetBlockData(0, data)

	c.Check(wrapped2.Get(res.ResourceID(1)).BlockData(0), check.DeepEquals, data)
}
