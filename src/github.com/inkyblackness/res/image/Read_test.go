package image

import (
	"bytes"
	"encoding/binary"

	check "gopkg.in/check.v1"
)

type ReadSuite struct {
}

var _ = check.Suite(&ReadSuite{})

func (suite *ReadSuite) TestReadReturnsErrorOnNilSource(c *check.C) {
	_, err := Read(nil)

	c.Check(err, check.ErrorMatches, "source is nil")
}

func (suite *ReadSuite) TestReadOfUncompressedDataReturnsBitmap(c *check.C) {
	data := suite.getTestData(UncompressedBitmap, []byte{0xAA}, false)
	bmp, err := Read(bytes.NewReader(data))

	c.Assert(err, check.IsNil)
	c.Assert(bmp, check.NotNil)
	c.Check(bmp.ImageWidth(), check.Equals, uint16(1))
	c.Check(bmp.ImageHeight(), check.Equals, uint16(1))
	c.Check(bmp.Row(0), check.DeepEquals, []byte{0xAA})
}

func (suite *ReadSuite) TestReadOfCompressedDataReturnsBitmap(c *check.C) {
	data := suite.getTestData(CompressedBitmap, []byte{0x00, 0x01, 0xBB, 0x80, 0x00, 0x00}, false)
	bmp, err := Read(bytes.NewReader(data))

	c.Assert(err, check.IsNil)
	c.Assert(bmp, check.NotNil)
	c.Check(bmp.ImageWidth(), check.Equals, uint16(1))
	c.Check(bmp.ImageHeight(), check.Equals, uint16(1))
	c.Check(bmp.Row(0), check.DeepEquals, []byte{0xBB})
}

func (suite *ReadSuite) TestReadWithPrivatePaletteReturnsPalette(c *check.C) {
	data := suite.getTestData(UncompressedBitmap, []byte{0xAA}, true)
	bmp, err := Read(bytes.NewReader(data))

	c.Assert(err, check.IsNil)
	c.Assert(bmp, check.NotNil)
	c.Check(bmp.Palette(), check.NotNil)
}

func (suite *ReadSuite) getTestData(bmpType BitmapType, data []byte, withPalette bool) []byte {
	var header BitmapHeader
	buf := bytes.NewBuffer(nil)

	header.Type = bmpType
	header.Width = 1
	header.Stride = 1
	header.Height = 1
	header.WidthFactor = 1
	header.HeightFactor = 1
	if withPalette {
		header.PaletteOffset = int32(binary.Size(header))
	}
	binary.Write(buf, binary.LittleEndian, &header)
	buf.Write(data)
	if withPalette {
		binary.Write(buf, binary.LittleEndian, privatePaletteFlag)
		buf.Write(make([]byte, ColorsPerPixel*bytesPerColor))
	}

	return buf.Bytes()
}
