package dos

import (
	"bytes"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/serial"

	check "gopkg.in/check.v1"
)

type FormatReaderSuite struct {
}

var _ = check.Suite(&FormatReaderSuite{})

func (suite *FormatReaderSuite) TestNewChunkProviderReturnsErrorOnNil(c *check.C) {
	_, err := NewChunkProvider(nil)

	c.Check(err, check.ErrorMatches, "source is nil")
}

func (suite *FormatReaderSuite) TestNewChunkProviderReturnsProviderOnEmptySource(c *check.C) {
	source := bytes.NewReader(emptyResourceFile())
	provider, _ := NewChunkProvider(source)

	c.Check(provider, check.NotNil)
}

func (suite *FormatReaderSuite) TestNewChunkProviderReturnsErrorOnInvalidHeaderString(c *check.C) {
	sourceData := emptyResourceFile()
	sourceData[10] = byte("A"[0])

	_, err := NewChunkProvider(bytes.NewReader(sourceData))

	c.Check(err, check.ErrorMatches, "Format mismatch")
}

func (suite *FormatReaderSuite) TestNewChunkProviderReturnsErrorOnMissingCommentTerminator(c *check.C) {
	sourceData := emptyResourceFile()
	sourceData[len(HeaderString)] = byte(0)

	_, err := NewChunkProvider(bytes.NewReader(sourceData))

	c.Check(err, check.ErrorMatches, "Format mismatch")
}

func (suite *FormatReaderSuite) TestNewChunkProviderReturnsErrorOnInvalidDirectoryStart(c *check.C) {
	sourceData := emptyResourceFile()
	sourceData[ChunkDirectoryFileOffsetPos] = byte(0xFF)

	_, err := NewChunkProvider(bytes.NewReader(sourceData))

	c.Check(err, check.ErrorMatches, "EOF")
}

func (suite *FormatReaderSuite) TestIDsReturnsTheStoredChunkIDsInOrder(c *check.C) {
	store := serial.NewByteStore()
	consumer := NewChunkConsumer(store)

	blockHolder1 := chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}})
	consumer.Consume(res.ResourceID(0x5678), blockHolder1)
	blockHolder2 := chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}})
	consumer.Consume(res.ResourceID(0x1234), blockHolder2)
	consumer.Finish()

	provider, _ := NewChunkProvider(bytes.NewReader(store.Data()))

	c.Check(provider.IDs(), check.DeepEquals, []res.ResourceID{0x5678, 0x1234})
}

func (suite *FormatReaderSuite) TestProvideReturnsABlockProviderForKnownID(c *check.C) {
	store := serial.NewByteStore()
	consumer := NewChunkConsumer(store)

	blockHolder1 := chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}})
	consumer.Consume(res.ResourceID(0x1122), blockHolder1)
	consumer.Finish()

	provider, _ := NewChunkProvider(bytes.NewReader(store.Data()))

	c.Check(provider.Provide(0x1122), check.NotNil)
}

func (suite *FormatReaderSuite) TestProvideReturnsABlockProviderWithContent(c *check.C) {
	store := serial.NewByteStore()
	consumer := NewChunkConsumer(store)

	blockHolder1 := chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{0xAA, 0xBB, 0xCC}})
	consumer.Consume(res.ResourceID(0x3344), blockHolder1)
	consumer.Finish()

	provider, _ := NewChunkProvider(bytes.NewReader(store.Data()))

	c.Check(provider.Provide(0x3344).BlockData(0), check.DeepEquals, []byte{0xAA, 0xBB, 0xCC})
}

func (suite *FormatReaderSuite) TestProvideReturnsABlockProviderWithMetaData(c *check.C) {
	store := serial.NewByteStore()
	consumer := NewChunkConsumer(store)

	blockHolder1 := chunk.NewBlockHolder(chunk.BasicChunkType, res.Bitmap, [][]byte{[]byte{0xAA, 0xBB, 0xCC}})
	consumer.Consume(res.ResourceID(0x3344), blockHolder1)
	consumer.Finish()

	provider, _ := NewChunkProvider(bytes.NewReader(store.Data()))
	holder := provider.Provide(0x3344)

	c.Check(holder.BlockCount(), check.Equals, uint16(1))
	c.Check(holder.ContentType(), check.Equals, res.Bitmap)
}

func (suite *FormatReaderSuite) TestProvideReturnsABlockProviderWithDictionaryContent(c *check.C) {
	store := serial.NewByteStore()
	consumer := NewChunkConsumer(store)

	blockHolder1 := chunk.NewBlockHolder(chunk.BasicChunkType.WithDirectory(), res.Palette,
		[][]byte{[]byte{0xAA, 0xBB, 0xCC}, []byte{0xDD, 0xEE, 0xFF}})
	consumer.Consume(res.ResourceID(0x3344), blockHolder1)
	consumer.Finish()

	provider, _ := NewChunkProvider(bytes.NewReader(store.Data()))
	holder := provider.Provide(0x3344)

	c.Check(holder.BlockCount(), check.Equals, uint16(2))
	c.Check(holder.BlockData(1), check.DeepEquals, []byte{0xDD, 0xEE, 0xFF})
}

func (suite *FormatReaderSuite) TestProvideReturnsABlockProviderWithCompressedDictionaryContent(c *check.C) {
	store := serial.NewByteStore()
	consumer := NewChunkConsumer(store)

	blockHolder1 := chunk.NewBlockHolder(chunk.BasicChunkType.WithDirectory().WithCompression(), res.Palette,
		[][]byte{[]byte{0x01, 0x02, 0x01, 0x02}, []byte{0x03, 0x04}, []byte{0x05}})
	consumer.Consume(res.ResourceID(0x4455), blockHolder1)
	consumer.Finish()

	provider, _ := NewChunkProvider(bytes.NewReader(store.Data()))

	c.Check(provider.Provide(0x4455).BlockData(2), check.DeepEquals, []byte{0x05})
}
