package io

import (
	"bytes"
	"github.com/inkyblackness/res/serial"
	"github.com/inkyblackness/res/textprop"
	"github.com/inkyblackness/res/textprop/dos"
	"github.com/inkyblackness/res/textprop/store"
)

import (
	check "gopkg.in/check.v1"
)

type DynamicTextPropStoreSuite struct {
}

var _ = check.Suite(&DynamicTextPropStoreSuite{})

func (suite *DynamicTextPropStoreSuite) SetUpTest(c *check.C) {
}

func (suite *DynamicTextPropStoreSuite) createProvider(filler func(consumer textprop.Consumer)) textprop.Provider {
	store := serial.NewByteStore()
	consumer := dos.NewConsumer(store)
	filler(consumer)
	consumer.Finish()

	provider, _ := dos.NewProvider(bytes.NewReader(store.Data()))

	return provider
}

func (suite *DynamicTextPropStoreSuite) testData(baseValue byte) []byte {
	data := make([]byte, textprop.TexturePropertiesLength)

	for i := 0; i < int(textprop.TexturePropertiesLength); i++ {
		data[i] = baseValue
	}

	return data
}

func (suite *DynamicTextPropStoreSuite) TestPutInsertsToWrapped(c *check.C) {
	provider := suite.createProvider(func(consumer textprop.Consumer) {
		consumer.Consume(0, suite.testData(1))
	})

	wrappedStore := store.NewProviderBacked(provider, func() {})
	testStore := NewDynamicTextPropStore(wrappedStore)
	newData := suite.testData(2)

	testStore.Put(2, newData)

	wrappedData := wrappedStore.Get(2)

	c.Check(wrappedData, check.DeepEquals, newData)
}

func (suite *DynamicTextPropStoreSuite) TestGetReturnsBlockFromWrapped(c *check.C) {
	initData := suite.testData(4)
	provider := suite.createProvider(func(consumer textprop.Consumer) {
		consumer.Consume(3, initData)
	})

	wrappedStore := store.NewProviderBacked(provider, func() {})
	testStore := NewDynamicTextPropStore(wrappedStore)

	retrievedData := testStore.Get(3)

	c.Check(retrievedData, check.DeepEquals, initData)
}

func (suite *DynamicTextPropStoreSuite) TestSwapReplacesWrapped(c *check.C) {
	secondData := suite.testData(6)
	provider0 := suite.createProvider(func(consumer textprop.Consumer) {
		consumer.Consume(10, suite.testData(3))
		consumer.Consume(11, suite.testData(4))
	})
	provider1 := suite.createProvider(func(consumer textprop.Consumer) {
		consumer.Consume(10, suite.testData(5))
		consumer.Consume(11, secondData)
	})

	testStore := NewDynamicTextPropStore(store.NewProviderBacked(provider0, func() {}))
	testStore.Swap(func(oldStore textprop.Store) textprop.Store {
		return store.NewProviderBacked(provider1, func() {})
	})

	retrievedData := testStore.Get(11)

	c.Check(retrievedData, check.DeepEquals, secondData)
}
