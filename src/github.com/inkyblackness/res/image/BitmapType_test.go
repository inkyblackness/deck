package image

import (
	check "gopkg.in/check.v1"
)

type BitmapTypeSuite struct {
}

var _ = check.Suite(&BitmapTypeSuite{})

func (suite *BitmapTypeSuite) TestString(c *check.C) {
	c.Check(UncompressedBitmap.String(), check.Equals, "Uncompressed")
	c.Check(CompressedBitmap.String(), check.Equals, "Compressed")
	c.Check(BitmapType(0x0123).String(), check.Equals, "Unknown (0x0123)")
}
