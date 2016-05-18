package base

import (
	"bytes"
	"io"

	"github.com/inkyblackness/res/serial"

	check "gopkg.in/check.v1"
)

type CompressorSuite struct {
	store      *serial.ByteStore
	compressor io.WriteCloser
}

var _ = check.Suite(&CompressorSuite{})

func (suite *CompressorSuite) SetUpTest(c *check.C) {
	suite.store = serial.NewByteStore()
	coder := serial.NewEncoder(suite.store)
	suite.compressor = NewCompressor(coder)
}

func (suite *CompressorSuite) TestWriteCompressesFirstReocurrence(c *check.C) {
	suite.compressor.Write([]byte{0x00, 0x01})
	suite.compressor.Write([]byte{0x00, 0x01})

	suite.compressor.Close()

	suite.thenWordsShouldBe(c, word(0x0000), word(0x0001), word(0x0100))
}

func (suite *CompressorSuite) TestWriteCompressesTest1(c *check.C) {
	suite.compressor.Write([]byte{0x00, 0x01, 0x00, 0x02, 0x01, 0x00, 0x01})
	suite.compressor.Close()

	suite.thenWordsShouldBe(c, word(0x0000), word(0x0001), word(0x0000), word(0x0002), word(0x0101), word(0x0001))
}

func (suite *CompressorSuite) thenWordsShouldBe(c *check.C, expected ...word) {
	source := bytes.NewReader(suite.store.Data())
	reader := newWordReader(serial.NewDecoder(source))
	var words []word

	for read := reader.read(); read != endOfStream; read = reader.read() {
		words = append(words, read)
	}

	c.Assert(words, check.DeepEquals, expected)
}
