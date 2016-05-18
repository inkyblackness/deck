package video

import (
	check "gopkg.in/check.v1"
)

type MaskstreamReaderSuite struct {
}

var _ = check.Suite(&MaskstreamReaderSuite{})

func (suite *MaskstreamReaderSuite) TestReadPanicsForMoreThan8Bytes(c *check.C) {
	reader := NewMaskstreamReader([]byte{})

	c.Check(func() { reader.Read(9) }, check.Panics, "Limit of byte count: 8")
}

func (suite *MaskstreamReaderSuite) TestReadPanicsForLessThan0Bytes(c *check.C) {
	reader := NewMaskstreamReader([]byte{})

	c.Check(func() { reader.Read(-1) }, check.Panics, "Minimum byte count: 0")
}

func (suite *MaskstreamReaderSuite) TestReadReturnsValueOfRequestedByteSize(c *check.C) {
	reader := NewMaskstreamReader([]byte{0xAF})

	result := reader.Read(1)

	c.Check(result, check.Equals, uint64(0xAF))
}

func (suite *MaskstreamReaderSuite) TestReadIntegerFromSourceInLittleEndianOrder(c *check.C) {
	reader := NewMaskstreamReader([]byte{0x11, 0x22})

	result := reader.Read(2)

	c.Check(result, check.Equals, uint64(0x2211))
}

func (suite *MaskstreamReaderSuite) TestReadCanProvideUpTo64BitValues(c *check.C) {
	reader := NewMaskstreamReader([]byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88})

	result := reader.Read(8)

	c.Check(result, check.Equals, uint64(0x8877665544332211))
}

func (suite *MaskstreamReaderSuite) TestReadFillsMissingBytesWithZero(c *check.C) {
	reader := NewMaskstreamReader([]byte{0xAA, 0xBB})

	result := reader.Read(8)

	c.Check(result, check.Equals, uint64(0x00BBAA))
}

func (suite *MaskstreamReaderSuite) TestReadAdvancesCurrentPosition(c *check.C) {
	reader := NewMaskstreamReader([]byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88})

	reader.Read(2)
	result := reader.Read(3)

	c.Check(result, check.Equals, uint64(0x554433))
}
