package objprop

import (
	"github.com/inkyblackness/res"

	check "gopkg.in/check.v1"
)

type ConstantsSuite struct {
}

var _ = check.Suite(&ConstantsSuite{})

func (suite *ConstantsSuite) TestStandardPropertiesReturnsProperLength(c *check.C) {
	descriptor := StandardProperties()
	totalLength := uint32(4)

	for _, classDesc := range descriptor {
		totalLength += classDesc.TotalDataLength()
	}

	c.Assert(totalLength, check.Equals, uint32(17951)) // as taken from original CD
}

func (suite *ConstantsSuite) TestStandardPropertiesReturnsProperAmount(c *check.C) {
	descriptor := StandardProperties()
	total := uint32(0)

	for _, classDesc := range descriptor {
		total += classDesc.TotalTypeCount()
	}

	c.Assert(total, check.Equals, uint32(476))
}

func (suite *ConstantsSuite) TestObjectIDToIndexReturnsNegative1ForUnknown(c *check.C) {
	result := ObjectIDToIndex(nil, res.MakeObjectID(0, 1, 2))

	c.Assert(result, check.Equals, -1)
}

func (suite *ConstantsSuite) TestObjectIDToIndexReturns0ForFirstEntry(c *check.C) {
	result := ObjectIDToIndex(StandardProperties(), res.MakeObjectID(0, 0, 0))

	c.Assert(result, check.Equals, 0)
}

func (suite *ConstantsSuite) TestObjectIDToIndexReturns475ForLastEntry(c *check.C) {
	result := ObjectIDToIndex(StandardProperties(), res.MakeObjectID(14, 4, 1))

	c.Assert(result, check.Equals, 475)
}
