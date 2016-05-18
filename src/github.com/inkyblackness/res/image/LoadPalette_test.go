package image

import (
	"bytes"
	"image/color"

	check "gopkg.in/check.v1"
)

type LoadPaletteSuite struct {
}

var _ = check.Suite(&LoadPaletteSuite{})

func (suite *LoadPaletteSuite) TestLoadPaletteReturnsErrorIfTooFewBytes(c *check.C) {
	reader := bytes.NewReader(make([]byte, ColorsPerPixel*bytesPerColor-1))
	_, err := LoadPalette(reader)

	c.Check(err, check.NotNil)
}

func (suite *LoadPaletteSuite) TestLoadPaletteReturnsPaletteObject(c *check.C) {
	reader := bytes.NewReader(make([]byte, ColorsPerPixel*bytesPerColor))
	palette, _ := LoadPalette(reader)

	c.Check(palette, check.NotNil)
}

func (suite *LoadPaletteSuite) TestLoadPaletteSetsFirstColorTransparent(c *check.C) {
	reader := bytes.NewReader(make([]byte, ColorsPerPixel*bytesPerColor))
	palette, _ := LoadPalette(reader)

	col := color.NRGBAModel.Convert(palette[0]).(color.NRGBA)

	c.Check(col.A, check.Equals, byte(0))
}
