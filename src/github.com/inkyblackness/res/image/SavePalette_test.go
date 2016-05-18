package image

import (
	"bytes"
	"image/color"

	check "gopkg.in/check.v1"
)

type SavePaletteSuite struct {
}

var _ = check.Suite(&SavePaletteSuite{})

func (suite *SavePaletteSuite) TestSavePaletteReturnsErrorIfInvalidLength(c *check.C) {
	palette := make([]color.Color, ColorsPerPixel-1)
	writer := bytes.NewBuffer(nil)

	err := SavePalette(writer, palette)

	c.Check(err, check.NotNil)
}

func (suite *SavePaletteSuite) TestSavePaletteWritesTriplets(c *check.C) {
	palette := suite.getTestPalette()
	writer := bytes.NewBuffer(nil)

	SavePalette(writer, palette)

	buf := writer.Bytes()
	c.Check(len(buf), check.Equals, ColorsPerPixel*bytesPerColor)
}

func (suite *SavePaletteSuite) TestSavePaletteWritesColor(c *check.C) {
	palette := suite.getTestPalette()
	writer := bytes.NewBuffer(nil)

	SavePalette(writer, palette)

	expected := make([]byte, ColorsPerPixel*bytesPerColor)
	offset := 0
	for i := 0; i < ColorsPerPixel; i++ {
		expected[offset+0] = 0xFF - byte(i)
		expected[offset+1] = 0xFF - byte(i)
		expected[offset+2] = 0xFF - byte(i)
		offset += 3
	}

	buf := writer.Bytes()
	c.Check(buf, check.DeepEquals, expected)
}

func (suite *SavePaletteSuite) getTestPalette() color.Palette {
	palette := make([]color.Color, ColorsPerPixel)

	for i := 0; i < ColorsPerPixel; i++ {
		alpha := byte(0xFF)

		if i == 0 {
			alpha = 0x00
		}
		palette[i] = color.NRGBA{R: 0xFF - byte(i), G: 0xFF - byte(i), B: 0xFF - byte(i), A: alpha}
	}

	return palette
}
