package io

import (
	"time"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	dosChunk "github.com/inkyblackness/res/chunk/dos"
	"github.com/inkyblackness/res/objprop"
	dosObjprop "github.com/inkyblackness/res/objprop/dos"
	"github.com/inkyblackness/res/textprop"
	dosTextprop "github.com/inkyblackness/res/textprop/dos"
	"github.com/inkyblackness/shocked-core/release"
)

import (
	check "gopkg.in/check.v1"
)

type ReleaseStoreLibrarySuite struct {
	source  release.Release
	sink    release.Release
	library StoreLibrary

	descriptors         []objprop.ClassDescriptor
	nullObjpropProvider objprop.Provider
}

var _ = check.Suite(&ReleaseStoreLibrarySuite{})

func (suite *ReleaseStoreLibrarySuite) SetUpTest(c *check.C) {
	suite.descriptors = objprop.StandardProperties()
	suite.nullObjpropProvider = objprop.NullProvider(suite.descriptors)

	suite.source = release.NewMemoryRelease()
	suite.sink = release.NewMemoryRelease()
	suite.library = NewReleaseStoreLibrary(suite.source, suite.sink, 0)
}

func (suite *ReleaseStoreLibrarySuite) createChunkResource(rel release.Release, name string, filler func(consumer chunk.Consumer)) {
	resource, _ := rel.NewResource(name, "")
	writer, _ := resource.AsSink()
	consumer := dosChunk.NewChunkConsumer(writer)
	filler(consumer)
	consumer.Finish()
}

func (suite *ReleaseStoreLibrarySuite) createObjpropResource(rel release.Release, name string, filler func(consumer objprop.Consumer)) {
	resource, _ := rel.NewResource(name, "")
	writer, _ := resource.AsSink()
	consumer := dosObjprop.NewConsumer(writer, suite.descriptors)
	filler(consumer)
	consumer.Finish()
}

func (suite *ReleaseStoreLibrarySuite) createTextpropResource(rel release.Release, name string, filler func(consumer textprop.Consumer)) {
	resource, _ := rel.NewResource(name, "")
	writer, _ := resource.AsSink()
	consumer := dosTextprop.NewConsumer(writer)
	filler(consumer)
	consumer.Finish()
}

func (suite *ReleaseStoreLibrarySuite) someObjectProperties(objID res.ObjectID) objprop.ObjectData {
	data := suite.nullObjpropProvider.Provide(objID)

	for index := 0; index < len(data.Common); index++ {
		data.Common[index] = byte(objID.Class)
	}
	for index := 0; index < len(data.Generic); index++ {
		data.Generic[index] = byte(objID.Subclass)
	}
	for index := 0; index < len(data.Specific); index++ {
		data.Specific[index] = byte(objID.Type)
	}

	return data
}

func (suite *ReleaseStoreLibrarySuite) someTextureProperties(seed byte) []byte {
	data := make([]byte, textprop.TexturePropertiesLength)

	for index := 0; index < len(data); index++ {
		data[index] = seed + byte(index)
	}

	return data
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

func (suite *ReleaseStoreLibrarySuite) TestModifyingChunkSourceSavesNewSink(c *check.C) {
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

func (suite *ReleaseStoreLibrarySuite) TestObjpropStoreIsBackedBySinkIfExisting(c *check.C) {
	objID := res.MakeObjectID(1, 2, 2)
	expected := suite.someObjectProperties(objID)
	suite.createObjpropResource(suite.sink, "objprop.dat", func(consumer objprop.Consumer) {
		consumer.Consume(objID, expected)
	})
	store, err := suite.library.ObjpropStore("objprop.dat")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	retrieved := store.Get(objID)
	c.Check(retrieved, check.DeepEquals, expected)
}

func (suite *ReleaseStoreLibrarySuite) TestObjpropStoreIsBackedBySinkIfExistingInBoth(c *check.C) {
	objID := res.MakeObjectID(1, 2, 2)
	expected := suite.someObjectProperties(objID)
	suite.createObjpropResource(suite.source, "stillFromSink.dat", func(consumer objprop.Consumer) {})
	suite.createObjpropResource(suite.sink, "stillFromSink.dat", func(consumer objprop.Consumer) {
		consumer.Consume(objID, expected)
	})
	store, err := suite.library.ObjpropStore("stillFromSink.dat")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	retrieved := store.Get(objID)
	c.Check(retrieved, check.DeepEquals, expected)
}

func (suite *ReleaseStoreLibrarySuite) TestObjpropStoreIsBackedBySourceIfExisting(c *check.C) {
	objID := res.MakeObjectID(1, 2, 2)
	expected := suite.someObjectProperties(objID)
	suite.createObjpropResource(suite.source, "fromSource.dat", func(consumer objprop.Consumer) {
		consumer.Consume(objID, expected)
	})
	store, err := suite.library.ObjpropStore("fromSource.dat")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	retrieved := store.Get(objID)
	c.Check(retrieved, check.DeepEquals, expected)
}

func (suite *ReleaseStoreLibrarySuite) TestObjpropStoreReturnsEmptyStoreIfNowhereExisting(c *check.C) {
	objID := res.MakeObjectID(1, 2, 2)
	expected := suite.nullObjpropProvider.Provide(objID)
	store, err := suite.library.ObjpropStore("empty.dat")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)

	retrieved := store.Get(objID)
	c.Check(retrieved, check.DeepEquals, expected)
}

func (suite *ReleaseStoreLibrarySuite) TestModifyingObjpropSourceSavesNewSink(c *check.C) {
	objID := res.MakeObjectID(1, 2, 2)
	nullData := suite.nullObjpropProvider.Provide(objID)
	suite.createObjpropResource(suite.source, "source.dat", func(consumer objprop.Consumer) {
		consumer.Consume(objID, suite.someObjectProperties(objID))
	})
	store, err := suite.library.ObjpropStore("source.dat")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)

	store.Put(objID, nullData)

	time.Sleep(100 * time.Millisecond)

	c.Check(suite.sink.HasResource("source.dat"), check.Equals, true)
}

func (suite *ReleaseStoreLibrarySuite) TestObjpropStoreReturnsSameInstances(c *check.C) {
	suite.createObjpropResource(suite.source, "source.dat", func(consumer objprop.Consumer) {
		objID := res.MakeObjectID(1, 2, 2)
		consumer.Consume(objID, suite.someObjectProperties(objID))
	})
	store1, err1 := suite.library.ObjpropStore("source.dat")
	c.Assert(err1, check.IsNil)
	c.Assert(store1, check.NotNil)

	store2, err2 := suite.library.ObjpropStore("source.dat")
	c.Assert(err2, check.IsNil)

	c.Check(store1, check.Equals, store2)
}

func (suite *ReleaseStoreLibrarySuite) TestTextpropStoreIsBackedBySinkIfExisting(c *check.C) {
	textID := uint32(123)
	expected := suite.someTextureProperties(0x10)
	suite.createTextpropResource(suite.sink, "textprop.dat", func(consumer textprop.Consumer) {
		consumer.Consume(textID, expected)
	})
	store, err := suite.library.TextpropStore("textprop.dat")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	retrieved := store.Get(textID)
	c.Check(retrieved, check.DeepEquals, expected)
}

func (suite *ReleaseStoreLibrarySuite) TestTextpropStoreIsBackedBySinkIfExistingInBoth(c *check.C) {
	textID := uint32(20)
	expected := suite.someTextureProperties(0x20)
	suite.createTextpropResource(suite.source, "stillFromSink.dat", func(consumer textprop.Consumer) {})
	suite.createTextpropResource(suite.sink, "stillFromSink.dat", func(consumer textprop.Consumer) {
		consumer.Consume(textID, expected)
	})
	store, err := suite.library.TextpropStore("stillFromSink.dat")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	retrieved := store.Get(textID)
	c.Check(retrieved, check.DeepEquals, expected)
}

func (suite *ReleaseStoreLibrarySuite) TestTextpropStoreIsBackedBySourceIfExisting(c *check.C) {
	textID := uint32(30)
	expected := suite.someTextureProperties(0x30)
	suite.createTextpropResource(suite.source, "fromSource.dat", func(consumer textprop.Consumer) {
		consumer.Consume(textID, expected)
	})
	store, err := suite.library.TextpropStore("fromSource.dat")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	retrieved := store.Get(textID)
	c.Check(retrieved, check.DeepEquals, expected)
}

func (suite *ReleaseStoreLibrarySuite) TestTextpropStoreReturnsEmptyStoreIfNowhereExisting(c *check.C) {
	textID := uint32(40)
	expected := make([]byte, textprop.TexturePropertiesLength)
	store, err := suite.library.TextpropStore("empty.dat")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)

	retrieved := store.Get(textID)
	c.Check(retrieved, check.DeepEquals, expected)
}

func (suite *ReleaseStoreLibrarySuite) TestModifyingTextpropSourceSavesNewSink(c *check.C) {
	textID := uint32(50)
	nullData := make([]byte, textprop.TexturePropertiesLength)
	suite.createTextpropResource(suite.source, "source.dat", func(consumer textprop.Consumer) {
		consumer.Consume(textID, suite.someTextureProperties(0x50))
	})
	store, err := suite.library.TextpropStore("source.dat")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)

	store.Put(textID, nullData)

	time.Sleep(100 * time.Millisecond)

	c.Check(suite.sink.HasResource("source.dat"), check.Equals, true)
}

func (suite *ReleaseStoreLibrarySuite) TestTextpropStoreReturnsSameInstances(c *check.C) {
	suite.createTextpropResource(suite.source, "source.dat", func(consumer textprop.Consumer) {
		textID := uint32(60)
		consumer.Consume(textID, suite.someTextureProperties(0x60))
	})
	store1, err1 := suite.library.TextpropStore("source.dat")
	c.Assert(err1, check.IsNil)
	c.Assert(store1, check.NotNil)

	store2, err2 := suite.library.TextpropStore("source.dat")
	c.Assert(err2, check.IsNil)

	c.Check(store1, check.Equals, store2)
}
