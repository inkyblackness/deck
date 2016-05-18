package serial

import (
	"bytes"

	check "gopkg.in/check.v1"
)

type PositioningDecoderSuite struct {
	coder PositioningCoder
}

var _ = check.Suite(&PositioningDecoderSuite{})

func (suite *PositioningDecoderSuite) SetUpTest(c *check.C) {
}

func (suite *PositioningDecoderSuite) TestSetCurPosRepositionsReadOffset(c *check.C) {
	var source = bytes.NewReader([]byte{0x78, 0x12, 0x34})
	arrayValue := make([]byte, 3)
	var intValue uint16

	suite.coder = NewPositioningDecoder(source)
	suite.coder.CodeBytes(arrayValue)
	suite.coder.SetCurPos(1)
	suite.coder.CodeUint16(&intValue)

	c.Assert(intValue, check.Equals, uint16(0x3412))
}
