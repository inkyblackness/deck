package objprop

import (
	check "gopkg.in/check.v1"
)

type ClassDescriptorSuite struct {
}

var _ = check.Suite(&ClassDescriptorSuite{})

func (suite *ClassDescriptorSuite) TestTotalDataLengthReturnsCompleteLength(c *check.C) {
	var mainDesc ClassDescriptor

	mainDesc.GenericDataLength = 7
	mainDesc.Subclasses = append(mainDesc.Subclasses, SubclassDescriptor{2, 3})
	mainDesc.Subclasses = append(mainDesc.Subclasses, SubclassDescriptor{1, 20})

	c.Assert(mainDesc.TotalDataLength(), check.Equals, uint32((7*(2+1))+(2*3)+(1*20)+CommonPropertiesLength*3))
}

func (suite *ClassDescriptorSuite) TestTotalTypeCountReturnsTotalAmount(c *check.C) {
	var mainDesc ClassDescriptor

	mainDesc.GenericDataLength = 7
	mainDesc.Subclasses = append(mainDesc.Subclasses, SubclassDescriptor{2, 3})
	mainDesc.Subclasses = append(mainDesc.Subclasses, SubclassDescriptor{1, 20})

	c.Assert(mainDesc.TotalTypeCount(), check.Equals, uint32(3))
}
