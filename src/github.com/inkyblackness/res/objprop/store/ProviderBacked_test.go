package store

import (
	"github.com/inkyblackness/res"

	check "gopkg.in/check.v1"
)

type ProviderBackedSuite struct {
}

var _ = check.Suite(&ProviderBackedSuite{})

func (suite *ProviderBackedSuite) TestGetReturnsPropertiesFromProvider_WhenUnmodified(c *check.C) {
	data1 := randomObjectData()
	id := res.MakeObjectID(1, 2, 3)
	provider := NewTestingProvider()
	provider.Consume(id, data1)
	backed := NewProviderBacked(provider, func() {})

	c.Check(backed.Get(id), check.DeepEquals, data1)
}

func (suite *ProviderBackedSuite) TestGetReturnsNewProperties_WhenModified(c *check.C) {
	data1 := randomObjectData()
	data2 := randomObjectData()
	id := res.MakeObjectID(1, 2, 3)
	provider := NewTestingProvider()
	provider.Consume(id, data1)
	backed := NewProviderBacked(provider, func() {})

	backed.Put(id, data2)

	c.Check(backed.Get(id), check.DeepEquals, data2)
}

func (suite *ProviderBackedSuite) TestModifiedCallback_WhenModified(c *check.C) {
	data1 := randomObjectData()
	data2 := randomObjectData()
	id := res.MakeObjectID(1, 2, 3)
	provider := NewTestingProvider()
	provider.Consume(id, data1)
	called := false
	backed := NewProviderBacked(provider, func() { called = true })

	backed.Put(id, data2)

	c.Check(called, check.Equals, true)
}
