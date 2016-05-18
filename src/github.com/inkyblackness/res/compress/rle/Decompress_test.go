package rle

import (
	"bytes"

	check "gopkg.in/check.v1"
)

type DecompressSuite struct {
}

var _ = check.Suite(&DecompressSuite{})

func (suite *DecompressSuite) TestEmptyArrayReturnsError(c *check.C) {
	err := Decompress(bytes.NewReader(nil), make([]byte, 100))

	c.Check(err, check.NotNil)
}

func (suite *DecompressSuite) Test800000IsEndOfStream(c *check.C) {
	result := []byte{}
	err := Decompress(bytes.NewReader([]byte{0x80, 0x00, 0x00}), result)

	c.Assert(err, check.IsNil)
	c.Check(result, check.DeepEquals, []byte{})
}

func (suite *DecompressSuite) Test800000AppendsRemainingZeroes(c *check.C) {
	result := make([]byte, 10)
	err := Decompress(bytes.NewReader([]byte{0x80, 0x00, 0x00}), result)

	c.Assert(err, check.IsNil)
	c.Check(result, check.DeepEquals, make([]byte, 10))
}

func (suite *DecompressSuite) Test800000IsConsumedAtEnd(c *check.C) {
	result := []byte{}
	reader := bytes.NewReader([]byte{0x80, 0x00, 0x00})
	err := Decompress(reader, result)

	c.Assert(err, check.IsNil)
	pos, _ := reader.Seek(0, 1)
	c.Check(pos, check.Equals, int64(3))
}

func (suite *DecompressSuite) Test00WritesNNBytesOfColorZZ(c *check.C) {
	result := make([]byte, 5)
	err := Decompress(bytes.NewReader([]byte{0x00, 0x05, 0xCC, 0x80, 0x00, 0x00}), result)

	c.Assert(err, check.IsNil)
	c.Check(result, check.DeepEquals, []byte{0xCC, 0xCC, 0xCC, 0xCC, 0xCC})
}

func (suite *DecompressSuite) Test00ReturnsErrorIfZZIsMissing(c *check.C) {
	err := Decompress(bytes.NewReader([]byte{0x00, 0x05}), make([]byte, 5))

	c.Check(err, check.NotNil)
}

func (suite *DecompressSuite) Test00ReturnsErrorIfNNIsMissing(c *check.C) {
	err := Decompress(bytes.NewReader([]byte{0x00}), make([]byte, 5))

	c.Check(err, check.NotNil)
}

func (suite *DecompressSuite) TestNNLess80WritesNNFollowingBytes(c *check.C) {
	result := make([]byte, 2)
	err := Decompress(bytes.NewReader([]byte{0x02, 0xAA, 0xBB, 0x80, 0x00, 0x00}), result)

	c.Assert(err, check.IsNil)
	c.Check(result, check.DeepEquals, []byte{0xAA, 0xBB})
}

func (suite *DecompressSuite) TestNNLess80ReturnsErrorIfEndOfFile(c *check.C) {
	err := Decompress(bytes.NewReader([]byte{0x02, 0xAA}), make([]byte, 2))

	c.Check(err, check.NotNil)
}

func (suite *DecompressSuite) Test80WritesZeroes(c *check.C) {
	result := make([]byte, 0x123)
	err := Decompress(bytes.NewReader([]byte{0x80, 0x23, 0x01, 0x80, 0x00, 0x00}), result)

	c.Assert(err, check.IsNil)
	c.Check(result, check.DeepEquals, make([]byte, 0x0123))
}

func (suite *DecompressSuite) Test80CopiesNNBytes(c *check.C) {
	input := bytes.NewBuffer(nil)
	input.Write([]byte{0x80, 0x04, 0x80})
	input.Write([]byte{0x01, 0x02, 0x03, 0x04})
	input.Write([]byte{0x80, 0x00, 0x00})
	result := make([]byte, 4)
	err := Decompress(bytes.NewReader(input.Bytes()), result)

	c.Assert(err, check.IsNil)
	c.Check(result, check.DeepEquals, []byte{0x01, 0x02, 0x03, 0x04})
}

func (suite *DecompressSuite) Test80CopiesNNBytesExtended(c *check.C) {
	input := bytes.NewBuffer(nil)
	expected := make([]byte, 0x3FFF)
	input.Write([]byte{0x80, 0xFF, 0xBF})
	input.Write(expected)
	input.Write([]byte{0x80, 0x00, 0x00})
	result := make([]byte, len(expected))
	err := Decompress(bytes.NewReader(input.Bytes()), result)

	c.Assert(err, check.IsNil)
	c.Check(result, check.DeepEquals, expected)
}

func (suite *DecompressSuite) Test80ReturnsErrorForUndefinedCase(c *check.C) {
	input := bytes.NewBuffer(nil)
	input.Write([]byte{0x80, 0x00, 0xC0})
	input.Write([]byte{0x80, 0x00, 0x00})
	err := Decompress(bytes.NewReader(input.Bytes()), make([]byte, 1))

	c.Check(err, check.NotNil)
}

func (suite *DecompressSuite) Test80WritesNNBytesOfValue(c *check.C) {
	input := bytes.NewBuffer(nil)
	expected := make([]byte, 0x3FFF)
	for index := range expected {
		expected[index] = 0xCD
	}
	input.Write([]byte{0x80, 0xFF, 0xFF, 0xCD})
	input.Write([]byte{0x80, 0x00, 0x00})
	result := make([]byte, len(expected))
	err := Decompress(bytes.NewReader(input.Bytes()), result)

	c.Assert(err, check.IsNil)
	c.Check(result, check.DeepEquals, expected)
}

func (suite *DecompressSuite) TestNNMore80WritesZeroes(c *check.C) {
	result := make([]byte, 3)
	err := Decompress(bytes.NewReader([]byte{0x83, 0x80, 0x00, 0x00}), result)

	c.Assert(err, check.IsNil)
	c.Check(result, check.DeepEquals, []byte{0, 0, 0})
}
