package image

import (
	"bytes"
	"encoding/binary"

	check "gopkg.in/check.v1"
)

type WriteSuite struct {
}

var _ = check.Suite(&WriteSuite{})

func (suite *WriteSuite) TestWriteUncompressedWithoutPalette(c *check.C) {
	sourceData := suite.getTestData(UncompressedBitmap, []byte{0xAA}, false)
	bmp, _ := Read(bytes.NewReader(sourceData))

	buf := bytes.NewBuffer(nil)
	Write(buf, bmp, UncompressedBitmap, 0)
	result := buf.Bytes()

	c.Check(result, check.DeepEquals, sourceData)
}

func (suite *WriteSuite) TestWriteCompressedWithoutPalette(c *check.C) {
	sourceData := suite.getTestData(CompressedBitmap, []byte{0x01, 0xBB, 0x80, 0x00, 0x00}, false)
	bmp, _ := Read(bytes.NewReader(sourceData))

	buf := bytes.NewBuffer(nil)
	Write(buf, bmp, CompressedBitmap, 0)
	result := buf.Bytes()

	c.Check(result, check.DeepEquals, sourceData)
}

func (suite *WriteSuite) TestWriteUncompressedWithPalette(c *check.C) {
	sourceData := suite.getTestData(UncompressedBitmap, []byte{0xCC}, true)
	bmp, _ := Read(bytes.NewReader(sourceData))

	buf := bytes.NewBuffer(nil)
	Write(buf, bmp, UncompressedBitmap, 0)
	result := buf.Bytes()

	c.Check(result, check.DeepEquals, sourceData)
}

func (suite *WriteSuite) TestWriteCompressedWithPalette(c *check.C) {
	sourceData := suite.getTestData(CompressedBitmap, []byte{0x01, 0xEE, 0x80, 0x00, 0x00}, true)
	bmp, _ := Read(bytes.NewReader(sourceData))

	buf := bytes.NewBuffer(nil)
	Write(buf, bmp, CompressedBitmap, 0)
	result := buf.Bytes()

	c.Check(result, check.DeepEquals, sourceData)
}

func (suite *WriteSuite) getTestData(bmpType BitmapType, data []byte, withPalette bool) []byte {
	var header BitmapHeader
	buf := bytes.NewBuffer(nil)

	header.Type = bmpType
	header.Width = 1
	header.Stride = 1
	header.Height = 1
	header.WidthFactor = 0
	header.HeightFactor = 0
	if withPalette {
		header.PaletteOffset = int32(binary.Size(header) + len(data))
	}
	binary.Write(buf, binary.LittleEndian, &header)
	buf.Write(data)
	if withPalette {
		binary.Write(buf, binary.LittleEndian, privatePaletteFlag)
		buf.Write(make([]byte, ColorsPerPixel*bytesPerColor))
	}

	return buf.Bytes()
}
