package movi

import (
	check "gopkg.in/check.v1"
)

type UtilSuite struct {
}

var _ = check.Suite(&UtilSuite{})

func (suite *UtilSuite) TestTimeFromRawA(c *check.C) {
	result := timeFromRaw(byte(4), uint16(0x8000))

	c.Check(result, check.Equals, float32(4.5))
}

func (suite *UtilSuite) TestTimeFromRawB(c *check.C) {
	result := timeFromRaw(byte(6), uint16(0))

	c.Check(result, check.Equals, float32(6.0))
}

func (suite *UtilSuite) TestTimeToRawA(c *check.C) {
	second, fraction := timeToRaw(7.25)

	c.Check(second, check.Equals, byte(7))
	c.Check(fraction, check.Equals, uint16(0x4000))
}

func (suite *UtilSuite) TestTimeToRawB(c *check.C) {
	second, fraction := timeToRaw(255.75)

	c.Check(second, check.Equals, byte(255))
	c.Check(fraction, check.Equals, uint16(0xC000))
}
