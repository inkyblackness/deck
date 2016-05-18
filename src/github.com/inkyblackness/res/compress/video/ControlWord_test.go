package video

import (
	check "gopkg.in/check.v1"
)

type ControlWordSuite struct {
}

var _ = check.Suite(&ControlWordSuite{})

func (suite *ControlWordSuite) TestPackedReturnsPackedControlWord(c *check.C) {
	word := ControlWord(0x00FFAA55)

	c.Check(word.Packed(128), check.Equals, PackedControlWord(0x80FFAA55))
}

func (suite *ControlWordSuite) TestPackedClearsHighestByteBeforeSettingCount(c *check.C) {
	word := ControlWord(0x55FFAA55)

	c.Check(word.Packed(0xC0), check.Equals, PackedControlWord(0xC0FFAA55))
}

func (suite *ControlWordSuite) TestCountReturnsBits20To23(c *check.C) {
	word := ControlWord(0x00F00000)

	c.Check(word.Count(), check.Equals, int(15))
}

func (suite *ControlWordSuite) TestIsLongOffsetReturnsTrueForCount0(c *check.C) {
	word := ControlWord(0x00000000)

	c.Check(word.IsLongOffset(), check.Equals, true)
}

func (suite *ControlWordSuite) TestIsLongOffsetReturnsFalseForCountGreater0(c *check.C) {
	word := ControlWord(0x00300000)

	c.Check(word.IsLongOffset(), check.Equals, false)
}

func (suite *ControlWordSuite) TestLongOffsetReturnsBits00To19(c *check.C) {
	word := ControlWord(0xFFFA6665)

	c.Check(word.LongOffset(), check.Equals, uint32(0xA6665))
}

func (suite *ControlWordSuite) TestParameterReturnsBits00To16(c *check.C) {
	word := ControlWord(0xFFFF1665)

	c.Check(word.Parameter(), check.Equals, uint32(0x11665))
}

func (suite *ControlWordSuite) TestTypeReturnsBits17To19(c *check.C) {
	word := ControlWord(0xFFFAFFFF)

	c.Check(word.Type(), check.Equals, ControlType(5))
}
