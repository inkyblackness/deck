package dos

import (
	"bytes"

	"github.com/inkyblackness/res/textprop"

	check "gopkg.in/check.v1"
)

type FormatReaderSuite struct {
}

var _ = check.Suite(&FormatReaderSuite{})

func (suite *FormatReaderSuite) TestNewProviderReturnsErrorOnNil(c *check.C) {
	_, err := NewProvider(nil)

	c.Assert(err, check.ErrorMatches, "source is nil")
}

func (suite *FormatReaderSuite) TestNewProviderReturnsErrorOnFileWithWrongHeader(c *check.C) {
	buf := bytes.NewBuffer(nil)
	buf.Write([]byte{0x08, 0x00, 0x00, 0x00})

	_, err := NewProvider(bytes.NewReader(buf.Bytes()))

	c.Assert(err, check.ErrorMatches, "Format mismatch")
}

func (suite *FormatReaderSuite) TestNewProviderReturnsProviderWithCountSet(c *check.C) {
	buf := bytes.NewBuffer(nil)
	buf.Write([]byte{0x09, 0x00, 0x00, 0x00})
	buf.Write(make([]byte, textprop.TexturePropertiesLength*4))
	source := bytes.NewReader(buf.Bytes())
	provider, err := NewProvider(source)

	c.Assert(err, check.IsNil)
	c.Check(provider.EntryCount(), check.Equals, uint32(4))
}

func (suite *FormatReaderSuite) TestNewProviderReturnsProviderWithData(c *check.C) {
	buf := bytes.NewBuffer(nil)
	buf.Write([]byte{0x09, 0x00, 0x00, 0x00})
	buf.Write(make([]byte, textprop.TexturePropertiesLength*2))
	sourceData := buf.Bytes()
	for i := byte(MagicHeaderSize); i < byte(len(sourceData)); i++ {
		sourceData[i] = i
	}
	source := bytes.NewReader(sourceData)
	provider, err := NewProvider(source)

	c.Assert(err, check.IsNil)
	c.Check(provider.Provide(1), check.DeepEquals, []byte{0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19})
}
