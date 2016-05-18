package font

import (
	"bytes"
	"encoding/binary"
	"io"

	check "gopkg.in/check.v1"
)

type LoadSuite struct {
}

var _ = check.Suite(&LoadSuite{})

func (suite *LoadSuite) SetUpTest(c *check.C) {
}

func (suite *LoadSuite) TestLoadReturnsErrorOnNil(c *check.C) {
	_, err := Load(nil)

	c.Check(err, check.ErrorMatches, "source is nil")
}

func (suite *LoadSuite) TestLoadReturnsFontForValidData(c *check.C) {
	font, err := Load(suite.aSimpleFont())

	c.Assert(err, check.IsNil)
	c.Check(font, check.NotNil)
}

func (suite *LoadSuite) TestLoadReturnsFontWithHeaderInfo(c *check.C) {
	font, err := Load(suite.aSimpleFont())

	c.Assert(err, check.IsNil)
	c.Check(font.IsMonochrome(), check.Equals, false)
	c.Check(font.FirstCharacter(), check.Equals, 32)
	c.Check(font.GlyphXOffset(0), check.Equals, 0)
	c.Check(font.GlyphXOffset(1), check.Equals, 1)
}

func (suite *LoadSuite) TestLoadReturnsFontWithBitmapData(c *check.C) {
	font, err := Load(suite.aSimpleFont())

	c.Assert(err, check.IsNil)
	c.Check(font.Bitmap(), check.DeepEquals, []byte{0xAB})
}

func (suite *LoadSuite) aSimpleFont() io.ReadSeeker {
	var header Header
	buf := bytes.NewBuffer(nil)

	header.Type = Color
	header.FirstCharacter = 32
	header.LastCharacter = 32
	header.XOffsetStart = uint32(HeaderSize)
	header.BitmapStart = header.XOffsetStart + uint32(4)
	header.Width = 1
	header.Height = 1
	binary.Write(buf, binary.LittleEndian, &header)
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(1))
	buf.Write([]byte{0xAB})

	return bytes.NewReader(buf.Bytes())
}
