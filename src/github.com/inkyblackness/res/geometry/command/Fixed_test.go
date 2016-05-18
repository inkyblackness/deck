package command

import (
	"fmt"

	check "gopkg.in/check.v1"
)

type FixedSuite struct {
}

var _ = check.Suite(&FixedSuite{})

func (suite *FixedSuite) TestToFixedWithZeroCreatesZeroValue(c *check.C) {
	coord := ToFixed(0.0)

	c.Check(uint32(coord), check.Equals, uint32(0))
}

func (suite *FixedSuite) TestStringPrintsValue(c *check.C) {
	coord := ToFixed(1.0)
	result := fmt.Sprintf("%v", coord)

	c.Check(result, check.Equals, "1")
}

func (suite *FixedSuite) TestFloatRetrievesValue_A(c *check.C) {
	coord := ToFixed(123.0)
	result := coord.Float()

	c.Check(result, check.Equals, float32(123.0))
}

func (suite *FixedSuite) TestFloatRetrievesValue_B(c *check.C) {
	coord := ToFixed(456.1211)
	result := coord.Float()

	c.Check(result, check.Equals, float32(456.1211))
}

func (suite *FixedSuite) TestFloatCanBeNegative(c *check.C) {
	coord := ToFixed(-123.5)
	result := coord.Float()

	c.Check(result, check.Equals, float32(-123.5))
}
