package store

import (
	check "gopkg.in/check.v1"
)

type ProviderBackedSuite struct {
}

var _ = check.Suite(&ProviderBackedSuite{})

func (suite *ProviderBackedSuite) TestEntryCountRetrievesValueFromProvider(c *check.C) {
	provider := NewTestingProvider()
	provider.Consume(uint32(1), randomProperties())
	provider.Consume(uint32(2), randomProperties())
	backed := NewProviderBacked(provider, func() {})

	c.Check(backed.EntryCount(), check.Equals, uint32(2))
}

func (suite *ProviderBackedSuite) TestGetReturnsPropertiesFromProvider_WhenUnmodified(c *check.C) {
	data1 := randomProperties()
	id := uint32(20)
	provider := NewTestingProvider()
	provider.Consume(id, data1)
	backed := NewProviderBacked(provider, func() {})

	c.Check(backed.Get(id), check.DeepEquals, data1)
}

func (suite *ProviderBackedSuite) TestGetReturnsNewProperties_WhenModified(c *check.C) {
	data1 := randomProperties()
	data2 := randomProperties()
	id := uint32(30)
	provider := NewTestingProvider()
	provider.Consume(id, data1)
	backed := NewProviderBacked(provider, func() {})

	backed.Put(id, data2)

	c.Check(backed.Get(id), check.DeepEquals, data2)
}

func (suite *ProviderBackedSuite) TestModifiedCallback_WhenModified(c *check.C) {
	data1 := randomProperties()
	data2 := randomProperties()
	id := uint32(40)
	provider := NewTestingProvider()
	provider.Consume(id, data1)
	called := false
	backed := NewProviderBacked(provider, func() { called = true })

	backed.Put(id, data2)

	c.Check(called, check.Equals, true)
}
