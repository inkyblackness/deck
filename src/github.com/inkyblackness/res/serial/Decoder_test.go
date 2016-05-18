package serial

import (
	"bytes"

	check "gopkg.in/check.v1"
)

type DecoderSuite struct {
	coder Coder
}

var _ = check.Suite(&DecoderSuite{})

func (suite *DecoderSuite) SetUpTest(c *check.C) {
}

func (suite *DecoderSuite) TestCodeUint32DecodesIntegerValue(c *check.C) {
	var source = bytes.NewReader([]byte{0x78, 0x56, 0x34, 0x12})
	var value uint32

	suite.coder = NewDecoder(source)
	suite.coder.CodeUint32(&value)

	c.Assert(value, check.Equals, uint32(0x12345678))
}

func (suite *DecoderSuite) TestCodeUint24DecodesIntegerValue(c *check.C) {
	var source = bytes.NewReader([]byte{0x78, 0x56, 0x34, 0x12})
	var value uint32

	suite.coder = NewDecoder(source)
	suite.coder.CodeUint24(&value)

	c.Assert(value, check.Equals, uint32(0x345678))
}

func (suite *DecoderSuite) TestCodeUint16DecodesIntegerValue(c *check.C) {
	var source = bytes.NewReader([]byte{0x45, 0x23})
	var value uint16

	suite.coder = NewDecoder(source)
	suite.coder.CodeUint16(&value)

	c.Assert(value, check.Equals, uint16(0x2345))
}

func (suite *DecoderSuite) TestCodeByteDecodesByteValue(c *check.C) {
	var source = bytes.NewReader([]byte{0x78})
	var value byte

	suite.coder = NewDecoder(source)
	suite.coder.CodeByte(&value)

	c.Assert(value, check.Equals, byte(0x78))
}

func (suite *DecoderSuite) TestCodeBytesDecodesByteArray(c *check.C) {
	var source = bytes.NewReader([]byte{0x78, 0x12, 0x34})
	value := make([]byte, 3)

	suite.coder = NewDecoder(source)
	suite.coder.CodeBytes(value)

	c.Assert(value, check.DeepEquals, []byte{0x78, 0x12, 0x34})
}
