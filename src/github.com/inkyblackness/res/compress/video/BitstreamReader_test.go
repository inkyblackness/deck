package video

import (
	check "gopkg.in/check.v1"
)

type BitstreamReaderSuite struct {
}

var _ = check.Suite(&BitstreamReaderSuite{})

func (suite *BitstreamReaderSuite) TestReadPanicsForMoreThan32Bits(c *check.C) {
	reader := NewBitstreamReader([]byte{0x11, 0x22, 0x33, 0x44})

	c.Check(func() { reader.Read(33) }, check.Panics, "Limit of bit count: 32")
}

func (suite *BitstreamReaderSuite) TestReadReturnsValueOfRequestedBitSize(c *check.C) {
	reader := NewBitstreamReader([]byte{0xAF})

	result := reader.Read(3)

	c.Check(result, check.Equals, uint32(5))
}

func (suite *BitstreamReaderSuite) TestRepeatedReadReturnsSameValue(c *check.C) {
	reader := NewBitstreamReader([]byte{0xAF})

	result1 := reader.Read(3)
	result2 := reader.Read(3)

	c.Check(result1, check.Equals, uint32(5))
	c.Check(result2, check.Equals, result1)
}

func (suite *BitstreamReaderSuite) TestReadReturnsZeroesForBitsBeyondEndA(c *check.C) {
	reader := NewBitstreamReader([]byte{0xAF})

	result := reader.Read(9)

	c.Check(result, check.Equals, uint32(0x15E))
}

func (suite *BitstreamReaderSuite) TestReadReturnsZeroesForBitsBeyondEndB(c *check.C) {
	reader := NewBitstreamReader([]byte{0xBF})

	result := reader.Read(32)

	c.Check(result, check.Equals, uint32(0xBF000000))
}

func (suite *BitstreamReaderSuite) TestAdvancePanicsForNegativeValues(c *check.C) {
	reader := NewBitstreamReader([]byte{0x11, 0x22, 0x33, 0x44})

	c.Check(func() { reader.Advance(-10) }, check.Panics, "Can only advance forward")
}

func (suite *BitstreamReaderSuite) TestAdvanceLetsReadFurtherBits(c *check.C) {
	reader := NewBitstreamReader([]byte{0xAF})

	reader.Advance(2)
	result := reader.Read(4)

	c.Check(result, check.Equals, uint32(0x0B))
}

func (suite *BitstreamReaderSuite) TestAdvanceToEndIsPossible(c *check.C) {
	reader := NewBitstreamReader([]byte{0xAF})

	reader.Advance(8)
	result := reader.Read(8)

	c.Check(result, check.Equals, uint32(0))
}

func (suite *BitstreamReaderSuite) TestInternalBufferDoesntLoseData(c *check.C) {
	reader := NewBitstreamReader([]byte{0x7F, 0xFF, 0xFF, 0xFF, 0x80})

	reader.Advance(1)
	result := reader.Read(32)

	c.Check(result, check.Equals, uint32(0xFFFFFFFF))
}

func (suite *BitstreamReaderSuite) TestAdvanceCanJumpToLastBit(c *check.C) {
	reader := NewBitstreamReader([]byte{0x00, 0x00, 0x01})

	reader.Advance(23)
	result := reader.Read(1)

	c.Check(result, check.Equals, uint32(1))
}

func (suite *BitstreamReaderSuite) TestReadAdvanceBeyondFirstRead(c *check.C) {
	reader := NewBitstreamReader([]byte{0xFF, 0x00, 0xFA})

	reader.Read(10)
	reader.Advance(20)
	result := reader.Read(4)

	c.Check(result, check.Equals, uint32(0x0A))
}

func (suite *BitstreamReaderSuite) TestReadAdvanceWithinFirstRead(c *check.C) {
	reader := NewBitstreamReader([]byte{0xFF, 0x00, 0xFA})

	reader.Read(10)
	reader.Advance(4)
	result := reader.Read(20)

	c.Check(result, check.Equals, uint32(0xF00FA))
}

func (suite *BitstreamReaderSuite) TestReadOfZeroBitsIsPossibleMidStream(c *check.C) {
	reader := NewBitstreamReader([]byte{0xFF})

	reader.Read(4)
	result := reader.Read(0)

	c.Check(result, check.Equals, uint32(0))
}

func (suite *BitstreamReaderSuite) TestReadOfZeroBitsIsPossibleAtEnd(c *check.C) {
	reader := NewBitstreamReader([]byte{0xFF})

	reader.Advance(8)
	result := reader.Read(0)

	c.Check(result, check.Equals, uint32(0))
}

func (suite *BitstreamReaderSuite) TestReadOfZeroBitsIsPossibleWithEmptySource(c *check.C) {
	reader := NewBitstreamReader([]byte{})

	result := reader.Read(0)

	c.Check(result, check.Equals, uint32(0))
}

func (suite *BitstreamReaderSuite) TestAdvanceBeyondEndIsPossible(c *check.C) {
	reader := NewBitstreamReader([]byte{0xFF})

	reader.Advance(9)
	result := reader.Read(32)

	c.Check(result, check.Equals, uint32(0))
}

func (suite *BitstreamReaderSuite) TestExhaustedReturnsFalseForAvailableBits(c *check.C) {
	reader := NewBitstreamReader([]byte{0xFF})

	result := reader.Exhausted()

	c.Check(result, check.Equals, false)
}

func (suite *BitstreamReaderSuite) TestExhaustedReturnsFalseForStillOneAvailableBit(c *check.C) {
	reader := NewBitstreamReader([]byte{0xFF})

	reader.Advance(7)
	result := reader.Exhausted()

	c.Check(result, check.Equals, false)
}

func (suite *BitstreamReaderSuite) TestExhaustedReturnsTrueAfterAdvancingToEnd(c *check.C) {
	reader := NewBitstreamReader([]byte{0xFF, 0xFF})

	reader.Advance(16)
	result := reader.Exhausted()

	c.Check(result, check.Equals, true)
}

func (suite *BitstreamReaderSuite) TestExhaustedReturnsTrueAfterAdvancingBeyondEnd(c *check.C) {
	reader := NewBitstreamReader([]byte{0xFF, 0xFF})

	reader.Advance(30)
	result := reader.Exhausted()

	c.Check(result, check.Equals, true)
}
