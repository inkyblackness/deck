package base

import (
	"github.com/inkyblackness/res/serial"

	check "gopkg.in/check.v1"
)

type WordWriterSuite struct {
	writer *wordWriter
	store  *serial.ByteStore
}

var _ = check.Suite(&WordWriterSuite{})

func (suite *WordWriterSuite) SetUpTest(c *check.C) {
	suite.store = serial.NewByteStore()
	coder := serial.NewEncoder(suite.store)
	suite.writer = newWordWriter(coder)
}

func (suite *WordWriterSuite) TestCloseWritesEndOfStreamMarkerAndTrailingZeroByte(c *check.C) {
	suite.writer.close()

	c.Assert(suite.store.Data(), check.DeepEquals, []byte{0xFF, 0xFC, 0x00})
}

func (suite *WordWriterSuite) TestCloseWritesRemainderOnlyIfNotEmpty(c *check.C) {
	suite.writer.write(word(0x0000))
	suite.writer.write(word(0x0000))
	suite.writer.write(word(0x0000))
	suite.writer.close()

	c.Assert(suite.store.Data(), check.DeepEquals, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x3F, 0xFF, 0x00})
}

func (suite *WordWriterSuite) TestWriteAndCloseLinesUpBits(c *check.C) {
	suite.writer.write(word(0x1FFE)) // 0111111 1111110
	suite.writer.close()             // 1111111 1111111

	c.Assert(suite.store.Data(), check.DeepEquals, []byte{0x7F, 0xFB, 0xFF, 0xF0, 0x00})
}
