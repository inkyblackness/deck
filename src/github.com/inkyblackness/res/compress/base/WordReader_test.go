package base

import (
	"bytes"

	"github.com/inkyblackness/res/serial"

	check "gopkg.in/check.v1"
)

type WordReaderSuite struct {
	reader *wordReader
}

var _ = check.Suite(&WordReaderSuite{})

func (suite *WordReaderSuite) givenReaderFrom(words ...word) {
	store := serial.NewByteStore()
	encoder := serial.NewEncoder(store)
	writer := newWordWriter(encoder)

	for _, value := range words {
		writer.write(value)
	}
	writer.close()

	source := bytes.NewReader(store.Data())
	suite.reader = newWordReader(serial.NewDecoder(source))
}

func (suite *WordReaderSuite) TestReadCanReadFirstWord(c *check.C) {
	suite.givenReaderFrom()

	c.Assert(suite.reader.read(), check.Equals, endOfStream)
}

func (suite *WordReaderSuite) TestReadCanReadSeveralWords(c *check.C) {
	input := []word{word(0x3FFF), word(0x0000), word(0x3FFF), word(0x0000), word(0x2001), word(0x1234)}
	suite.givenReaderFrom(input...)

	output := make([]word, len(input)+1)
	for i := 0; i < len(output); i++ {
		output[i] = suite.reader.read()
	}

	c.Assert(output, check.DeepEquals, append(input, endOfStream))
}
