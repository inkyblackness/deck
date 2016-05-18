package rle

import (
	"bytes"

	check "gopkg.in/check.v1"
)

type CompressSuite struct {
}

var _ = check.Suite(&CompressSuite{})

func (suite *CompressSuite) TestEmptyArrayResultsInTerminator(c *check.C) {
	writer := bytes.NewBuffer(nil)

	Compress(writer, nil)

	c.Check(writer.Bytes(), check.DeepEquals, []byte{0x80, 0x00, 0x00})
}

func (suite *CompressSuite) TestWriteZeroOfLessThan80(c *check.C) {
	writer := bytes.NewBuffer(nil)

	writeZero(writer, 0x7F)

	c.Check(writer.Bytes(), check.DeepEquals, []byte{0xFF})
}

func (suite *CompressSuite) TestWriteZeroOfLessThanFF(c *check.C) {
	writer := bytes.NewBuffer(nil)

	writeZero(writer, 0xFD)

	c.Check(writer.Bytes(), check.DeepEquals, []byte{0xFF, 0xFE})
}

func (suite *CompressSuite) TestWriteZeroOfLessThan8000(c *check.C) {
	writer := bytes.NewBuffer(nil)

	writeZero(writer, 0x7FFC)

	c.Check(writer.Bytes(), check.DeepEquals, []byte{0x80, 0xFC, 0x7F})
}

func (suite *CompressSuite) TestWriteZeroOf8000(c *check.C) {
	writer := bytes.NewBuffer(nil)

	writeZero(writer, 0x8000)

	c.Check(writer.Bytes(), check.DeepEquals, []byte{0x80, 0xFF, 0x7F, 0x81})
}

func (suite *CompressSuite) TestWriteRawOfLessThan80(c *check.C) {
	writer := bytes.NewBuffer(nil)

	writeRaw(writer, []byte{0x0A, 0x0B, 0x0C})

	c.Check(writer.Bytes(), check.DeepEquals, []byte{0x03, 0x0A, 0x0B, 0x0C})
}
