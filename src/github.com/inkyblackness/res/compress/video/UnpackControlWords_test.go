package video

import (
	check "gopkg.in/check.v1"
)

type UnpackControlWordsSuite struct {
}

var _ = check.Suite(&UnpackControlWordsSuite{})

func (suite *UnpackControlWordsSuite) TestFormatErrorForEmptyArray(c *check.C) {
	_, err := UnpackControlWords(nil)

	c.Check(err, check.Equals, FormatError)
}

func (suite *UnpackControlWordsSuite) TestFormatErrorForTooSmallSizeField(c *check.C) {
	_, err := UnpackControlWords(make([]byte, 3))

	c.Check(err, check.Equals, FormatError)
}

func (suite *UnpackControlWordsSuite) TestFormatErrorForSizeValueNotMultipleOfThree(c *check.C) {
	_, err := UnpackControlWords([]byte{0x02, 0x00, 0x00, 0x00, 0xCC, 0xBB, 0xAA, 0x01})

	c.Check(err, check.Equals, FormatError)
}

func (suite *UnpackControlWordsSuite) TestEmptyResultForNoEntries(c *check.C) {
	words, err := UnpackControlWords([]byte{0x00, 0x00, 0x00, 0x00})

	c.Assert(err, check.IsNil)
	c.Check(len(words), check.Equals, 0)
}

func (suite *UnpackControlWordsSuite) TestSingleEntry(c *check.C) {
	words, err := UnpackControlWords([]byte{0x03, 0x00, 0x00, 0x00, 0xCC, 0xBB, 0xAA, 0x01})

	c.Assert(err, check.IsNil)
	c.Check(words, check.DeepEquals, []ControlWord{ControlWord(0xAABBCC)})
}

func (suite *UnpackControlWordsSuite) TestMultipleEntry(c *check.C) {
	words, err := UnpackControlWords([]byte{0x09, 0x00, 0x00, 0x00, 0xCC, 0xBB, 0xAA, 0x02, 0x33, 0x22, 0x11, 0x01})

	c.Assert(err, check.IsNil)
	c.Check(words, check.DeepEquals, []ControlWord{ControlWord(0xAABBCC), ControlWord(0xAABBCC), ControlWord(0x112233)})
}

func (suite *UnpackControlWordsSuite) TestErrorForTooManyUnpacked(c *check.C) {
	_, err := UnpackControlWords([]byte{0x03, 0x00, 0x00, 0x00, 0xCC, 0xBB, 0xAA, 0x02})

	c.Check(err, check.NotNil)
}
