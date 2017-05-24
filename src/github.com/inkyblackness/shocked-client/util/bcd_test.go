package util

import (
	check "gopkg.in/check.v1"
)

type BCDSuite struct {
}

var _ = check.Suite(&BCDSuite{})

func (suite *BCDSuite) TestToBinaryCodedDecimal(c *check.C) {

	c.Check(ToBinaryCodedDecimal(0), check.Equals, uint16(0x0000))
	c.Check(ToBinaryCodedDecimal(1), check.Equals, uint16(0x0001))
	c.Check(ToBinaryCodedDecimal(20), check.Equals, uint16(0x0020))
	c.Check(ToBinaryCodedDecimal(300), check.Equals, uint16(0x0300))
	c.Check(ToBinaryCodedDecimal(4000), check.Equals, uint16(0x4000))
	c.Check(ToBinaryCodedDecimal(5678), check.Equals, uint16(0x5678))
	c.Check(ToBinaryCodedDecimal(23456), check.Equals, uint16(0x3456))
}

func (suite *BCDSuite) TestFromBinaryCodedDecimal(c *check.C) {
	c.Check(FromBinaryCodedDecimal(0x0000), check.Equals, uint16(0))
	c.Check(FromBinaryCodedDecimal(0x0001), check.Equals, uint16(1))
	c.Check(FromBinaryCodedDecimal(0x0020), check.Equals, uint16(20))
	c.Check(FromBinaryCodedDecimal(0x0300), check.Equals, uint16(300))
	c.Check(FromBinaryCodedDecimal(0x4000), check.Equals, uint16(4000))
	c.Check(FromBinaryCodedDecimal(0x5678), check.Equals, uint16(5678))
	c.Check(FromBinaryCodedDecimal(0x3456), check.Equals, uint16(3456))
}
