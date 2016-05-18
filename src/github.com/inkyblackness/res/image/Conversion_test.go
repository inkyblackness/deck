package image

import (
	"bytes"
	"encoding/binary"
	//	"image"

	check "gopkg.in/check.v1"
)

type ConversionSuite struct {
}

var _ = check.Suite(&ConversionSuite{})

func (suite *ConversionSuite) TestFromAndToWithoutPrivatePalette(c *check.C) {
	sourceData := suite.getTestData(UncompressedBitmap, []byte{0xAA}, false)
	sourceBitmap, _ := Read(bytes.NewReader(sourceData))
	img := FromBitmap(sourceBitmap, nil)
	destBitmap := ToBitmap(img, nil)

	buf := bytes.NewBuffer(nil)

	binary.Write(buf, binary.LittleEndian, &destBitmap.header)
	binary.Write(buf, binary.LittleEndian, destBitmap.data)

	destData := buf.Bytes()

	c.Check(destData, check.DeepEquals, sourceData)
}

func (suite *ConversionSuite) TestFromAndToWithPrivatePalette(c *check.C) {
	sourceData := suite.getTestData(UncompressedBitmap, []byte{0xBB}, true)
	sourceBitmap, _ := Read(bytes.NewReader(sourceData))
	img := FromBitmap(sourceBitmap, nil)
	destBitmap := ToBitmap(img, sourceBitmap.Palette())

	buf := bytes.NewBuffer(nil)

	destBitmap.header.PaletteOffset = int32(binary.Size(destBitmap.header))
	binary.Write(buf, binary.LittleEndian, &destBitmap.header)
	binary.Write(buf, binary.LittleEndian, destBitmap.data)
	binary.Write(buf, binary.LittleEndian, privatePaletteFlag)
	SavePalette(buf, destBitmap.Palette())

	destData := buf.Bytes()
	c.Check(destData, check.DeepEquals, sourceData)
}

func (suite *ConversionSuite) getTestData(bmpType BitmapType, data []byte, withPalette bool) []byte {
	var header BitmapHeader
	buf := bytes.NewBuffer(nil)

	header.Type = bmpType
	header.Width = 1
	header.Stride = 1
	header.Height = 1
	header.WidthFactor = 0
	header.HeightFactor = 0
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
