package video

import (
	check "gopkg.in/check.v1"
)

type PackedControlWordSuite struct {
}

var _ = check.Suite(&PackedControlWordSuite{})

func (suite *PackedControlWordSuite) TestTimesReturnsHighestByte(c *check.C) {
	packed := PackedControlWord(0xFF706050)

	c.Check(packed.Times(), check.Equals, int(0xFF))
}

func (suite *PackedControlWordSuite) TestValueReturnsTheControlWord(c *check.C) {
	packed := PackedControlWord(0xFF706050)

	c.Check(uint32(packed.Value()), check.Equals, uint32(0x00706050))
}
