package dos

import (
	//"bytes"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/serial"

	check "gopkg.in/check.v1"
)

type FormatWriterSuite struct {
	store    *serial.ByteStore
	consumer chunk.Consumer
}

var _ = check.Suite(&FormatWriterSuite{})

func (suite *FormatWriterSuite) SetUpTest(c *check.C) {
	suite.store = serial.NewByteStore()
	suite.consumer = NewChunkConsumer(suite.store)
}

func (suite *FormatWriterSuite) TestFinishWithoutAddingCreatesValidFileWithoutChunks(c *check.C) {
	expected := emptyResourceFile()

	suite.consumer.Finish()
	result := suite.store.Data()

	c.Assert(result, check.DeepEquals, expected)
}

func (suite *FormatWriterSuite) TestConsumeOfFlatUncompressedChunkCanBeWritten(c *check.C) {
	singleBlock := []byte{0xAB, 0x01, 0xCD, 0x02, 0xEF}
	blockHolder := chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{singleBlock})

	suite.consumer.Consume(res.ResourceID(0x1234), blockHolder)
	suite.consumer.Finish()

	result := suite.store.Data()

	expected := []byte{}
	expected = append(expected, singleBlock...)
	expected = append(expected, 0x00, 0x00, 0x00)       // alignment for directory
	expected = append(expected, 0x01, 0x00)             // chunk count
	expected = append(expected, 0x80, 0x00, 0x00, 0x00) // offset to first chunk
	expected = append(expected, 0x34, 0x12)             // chunk ID
	expected = append(expected, 0x05, 0x00, 0x00)       // chunk length (uncompressed)
	expected = append(expected, 0x00)                   // chunk type
	expected = append(expected, 0x05, 0x00, 0x00)       // chunk length in file
	expected = append(expected, 0x00)                   // content type
	c.Assert(result[ChunkDirectoryFileOffsetPos+4:], check.DeepEquals, expected)
}

func (suite *FormatWriterSuite) TestConsumeOfDirUncompressedChunkCanBeWritten(c *check.C) {
	singleBlock1 := []byte{0xAB, 0x01, 0xCD}
	singleBlock2 := []byte{0x11, 0x22, 0x33, 0x44}
	blockHolder := chunk.NewBlockHolder(chunk.BasicChunkType.WithDirectory(), res.Palette, [][]byte{singleBlock1, singleBlock2})

	suite.consumer.Consume(res.ResourceID(0x5678), blockHolder)
	suite.consumer.Finish()

	result := suite.store.Data()

	expected := []byte{}

	expected = append(expected, 0x02, 0x00)             // number of blocks
	expected = append(expected, 0x0E, 0x00, 0x00, 0x00) // offset to first block
	expected = append(expected, 0x11, 0x00, 0x00, 0x00) // offset to second block
	expected = append(expected, 0x15, 0x00, 0x00, 0x00) // size of chunk
	expected = append(expected, singleBlock1...)
	expected = append(expected, singleBlock2...)
	expected = append(expected, 0x00, 0x00, 0x00)       // alignment for directory
	expected = append(expected, 0x01, 0x00)             // chunk count
	expected = append(expected, 0x80, 0x00, 0x00, 0x00) // offset to first chunk
	expected = append(expected, 0x78, 0x56)             // chunk ID
	expected = append(expected, 0x15, 0x00, 0x00)       // chunk length (uncompressed)
	expected = append(expected, 0x02)                   // chunk type
	expected = append(expected, 0x15, 0x00, 0x00)       // chunk length in file
	expected = append(expected, 0x00)                   // content type
	c.Assert(result[ChunkDirectoryFileOffsetPos+4:], check.DeepEquals, expected)
}

func (suite *FormatWriterSuite) TestConsumeOfFlatCompressedChunkCanBeWritten(c *check.C) {
	singleBlock := []byte{0x01, 0x02, 0x01, 0x02}
	blockHolder := chunk.NewBlockHolder(chunk.BasicChunkType.WithCompression(), res.Palette, [][]byte{singleBlock})

	suite.consumer.Consume(res.ResourceID(0x1122), blockHolder)
	suite.consumer.Finish()

	result := suite.store.Data()

	expected := []byte{}
	// 0000 0000|0000 0100|0000 0000|0010 0000|0100 0000|0011 1111|1111 1111
	expected = append(expected, 0x00, 0x04, 0x00, 0x20, 0x40, 0x3F, 0xFF, 0x00) // 14bit words 0x0001 0x0002 0x0100 0x3FFF + trailing 0x00
	expected = append(expected)                                                 // alignment for directory
	expected = append(expected, 0x01, 0x00)                                     // chunk count
	expected = append(expected, 0x80, 0x00, 0x00, 0x00)                         // offset to first chunk
	expected = append(expected, 0x22, 0x11)                                     // chunk ID
	expected = append(expected, 0x04, 0x00, 0x00)                               // chunk length (uncompressed)
	expected = append(expected, 0x01)                                           // chunk type
	expected = append(expected, 0x08, 0x00, 0x00)                               // chunk length in file
	expected = append(expected, 0x00)                                           // content type
	c.Assert(result[ChunkDirectoryFileOffsetPos+4:], check.DeepEquals, expected)
}
