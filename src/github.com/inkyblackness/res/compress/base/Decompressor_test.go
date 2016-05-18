package base

import (
	"bytes"
	"io"
	"math/rand"
	"time"

	"github.com/inkyblackness/res/serial"

	check "gopkg.in/check.v1"
)

type DecompressorSuite struct {
	store      *serial.ByteStore
	compressor io.WriteCloser
}

var _ = check.Suite(&DecompressorSuite{})

func (suite *DecompressorSuite) SetUpTest(c *check.C) {
}

func (suite *DecompressorSuite) TestDecompressTest1(c *check.C) {
	input := []byte{0x00, 0x01, 0x00, 0x01}

	suite.verify(c, input)
}

func (suite *DecompressorSuite) TestDecompressTest2(c *check.C) {
	input := []byte{0x00, 0x01, 0x00, 0x01, 0x00, 0x01}

	suite.verify(c, input)
}

func (suite *DecompressorSuite) TestDecompressTest3(c *check.C) {
	input := []byte{}

	suite.verify(c, input)
}

func (suite *DecompressorSuite) TestDecompressTest4(c *check.C) {
	input := []byte{0x00, 0x01, 0x00, 0x02, 0x01, 0x00, 0x01}

	suite.verify(c, input)
}

func (suite *DecompressorSuite) TestDecompressTest5(c *check.C) {
	input := []byte{0x00, 0x01, 0x00, 0x02, 0x01, 0x00, 0x01, 0x02, 0x01, 0x02, 0x01, 0x00, 0x01, 0x02}

	suite.verify(c, input)
}

func (suite *DecompressorSuite) TestDecompressTestRandom(c *check.C) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for testCase := 0; testCase < 100; testCase++ {
		input := make([]byte, r.Intn(1024))
		for i := 0; i < len(input); i++ {
			input[i] = byte(r.Intn(256))
		}
		suite.verify(c, input)
	}
}

func (suite *DecompressorSuite) TestDecompressHandlesDictionaryResets(c *check.C) {
	suite.writeWords(0x0001, 0x0002, 0x0100, reset, 0x0003, 0x0004, 0x0100, endOfStream)

	suite.verifyOutput(c, []byte{0x01, 0x02, 0x01, 0x02, 0x03, 0x04, 0x03, 0x04})
}

func (suite *DecompressorSuite) TestDecompressHandlesSelfReferencingWords(c *check.C) {
	suite.writeWords(0x0001, 0x0002, 0x0101, endOfStream)

	suite.verifyOutput(c, []byte{0x01, 0x02, 0x02, 0x02})
}

func (suite *DecompressorSuite) writeWords(values ...word) {
	suite.store = serial.NewByteStore()
	coder := serial.NewEncoder(suite.store)
	writer := newWordWriter(coder)

	for _, value := range values {
		writer.write(value)
	}
	writer.close()
}

func (suite *DecompressorSuite) verify(c *check.C, input []byte) {
	suite.store = serial.NewByteStore()
	suite.compressor = NewCompressor(serial.NewEncoder(suite.store))

	suite.compressor.Write(input)
	suite.compressor.Close()

	suite.verifyOutput(c, input)
}

func (suite *DecompressorSuite) verifyOutput(c *check.C, expected []byte) {
	output := suite.buffer(len(expected))
	source := bytes.NewReader(suite.store.Data())
	decompressor := NewDecompressor(serial.NewDecoder(source))
	decompressor.Read(output)

	c.Check(output, check.DeepEquals, expected)
}

func (suite *DecompressorSuite) buffer(byteCount int) []byte {
	result := make([]byte, byteCount)
	for i := range result {
		result[i] = 0xFF
	}
	return result
}
