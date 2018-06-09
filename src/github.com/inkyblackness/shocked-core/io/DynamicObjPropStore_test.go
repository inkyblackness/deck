package io

import (
	"bytes"
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/objprop"
	"github.com/inkyblackness/res/objprop/dos"
	"github.com/inkyblackness/res/objprop/store"
	"github.com/inkyblackness/res/serial"
)

import (
	check "gopkg.in/check.v1"
)

type DynamicObjPropStoreSuite struct {
	descriptors []objprop.ClassDescriptor
}

var _ = check.Suite(&DynamicObjPropStoreSuite{})

func (suite *DynamicObjPropStoreSuite) SetUpTest(c *check.C) {
	suite.descriptors = []objprop.ClassDescriptor{}

	subclasses := []objprop.SubclassDescriptor{
		objprop.SubclassDescriptor{TypeCount: 2, SpecificDataLength: 1}}

	suite.descriptors = append(suite.descriptors, objprop.ClassDescriptor{GenericDataLength: 2, Subclasses: subclasses})
}

func (suite *DynamicObjPropStoreSuite) createProvider(filler func(consumer objprop.Consumer)) objprop.Provider {
	store := serial.NewByteStore()
	consumer := dos.NewConsumer(store, suite.descriptors)
	filler(consumer)
	consumer.Finish()

	provider, _ := dos.NewProvider(bytes.NewReader(store.Data()), suite.descriptors)

	return provider
}

func (suite *DynamicObjPropStoreSuite) testData(baseValue byte) objprop.ObjectData {
	var data objprop.ObjectData

	data.Common = make([]byte, objprop.CommonPropertiesLength)
	for i := 0; i < int(objprop.CommonPropertiesLength); i++ {
		data.Common[i] = baseValue
	}
	data.Generic = []byte{0x10 + baseValue, 0x10 + baseValue}
	data.Specific = []byte{0x20 + baseValue}

	return data
}

func (suite *DynamicObjPropStoreSuite) TestPutInsertsToWrapped(c *check.C) {
	objID := res.MakeObjectID(res.ObjectClass(0), res.ObjectSubclass(0), res.ObjectType(0))
	provider := suite.createProvider(func(consumer objprop.Consumer) {
		consumer.Consume(objID, suite.testData(0))
	})

	wrappedStore := store.NewProviderBacked(provider, func() {})
	testStore := NewDynamicObjPropStore(wrappedStore)
	newData := suite.testData(1)

	testStore.Put(objID, newData)

	wrappedData := wrappedStore.Get(objID)

	c.Check(wrappedData.Common, check.DeepEquals, newData.Common)
	c.Check(wrappedData.Generic, check.DeepEquals, newData.Generic)
	c.Check(wrappedData.Specific, check.DeepEquals, newData.Specific)
}

func (suite *DynamicObjPropStoreSuite) TestGetReturnsBlockFromWrapped(c *check.C) {
	objID := res.MakeObjectID(res.ObjectClass(0), res.ObjectSubclass(0), res.ObjectType(0))
	initData := suite.testData(4)
	provider := suite.createProvider(func(consumer objprop.Consumer) {
		consumer.Consume(objID, initData)
	})

	wrappedStore := store.NewProviderBacked(provider, func() {})
	testStore := NewDynamicObjPropStore(wrappedStore)

	retrievedData := testStore.Get(objID)

	c.Check(retrievedData.Common, check.DeepEquals, initData.Common)
	c.Check(retrievedData.Generic, check.DeepEquals, initData.Generic)
	c.Check(retrievedData.Specific, check.DeepEquals, initData.Specific)
}

func (suite *DynamicObjPropStoreSuite) TestSwapReplacesWrapped(c *check.C) {
	objID0 := res.MakeObjectID(res.ObjectClass(0), res.ObjectSubclass(0), res.ObjectType(0))
	objID1 := res.MakeObjectID(res.ObjectClass(0), res.ObjectSubclass(0), res.ObjectType(1))
	secondData := suite.testData(6)
	provider0 := suite.createProvider(func(consumer objprop.Consumer) {
		consumer.Consume(objID0, suite.testData(3))
		consumer.Consume(objID1, suite.testData(3))
	})
	provider1 := suite.createProvider(func(consumer objprop.Consumer) {
		consumer.Consume(objID0, suite.testData(5))
		consumer.Consume(objID1, secondData)
	})

	testStore := NewDynamicObjPropStore(store.NewProviderBacked(provider0, func() {}))
	testStore.Swap(func(oldStore objprop.Store) objprop.Store {
		return store.NewProviderBacked(provider1, func() {})
	})

	retrievedData := testStore.Get(objID1)

	c.Check(retrievedData.Common, check.DeepEquals, secondData.Common)
	c.Check(retrievedData.Generic, check.DeepEquals, secondData.Generic)
	c.Check(retrievedData.Specific, check.DeepEquals, secondData.Specific)
}
