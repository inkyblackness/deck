package movi

import (
	"bytes"

	check "gopkg.in/check.v1"

	"github.com/inkyblackness/res/movi/format"
)

type ReadSuite struct {
}

var _ = check.Suite(&ReadSuite{})

func (suite *ReadSuite) TestReadReturnsErrorOnNil(c *check.C) {
	_, err := Read(nil)

	c.Check(err, check.ErrorMatches, "source is nil")
}

func (suite *ReadSuite) TestReadReturnsContainerOnEmptyFile(c *check.C) {
	buffer := bytes.NewBufferString(format.Tag)
	buffer.Write(make([]byte, 0x100+0x300-len(format.Tag)))
	emptyFile := buffer.Bytes()
	source := bytes.NewReader(emptyFile)
	container, _ := Read(source)

	c.Check(container, check.NotNil)
}

func (suite *ReadSuite) TestReadReturnsErrorOnMissingTag(c *check.C) {
	emptyFile := make([]byte, 0x100+0x300)
	source := bytes.NewReader(emptyFile)
	_, err := Read(source)

	c.Check(err, check.ErrorMatches, "Not a MOVI format")
}

func (suite *ReadSuite) TestReadReturnsContainerWithBasicPropertiesSet(c *check.C) {
	buffer := bytes.NewBufferString(format.Tag)
	buffer.Write(make([]byte, 0x100+0x300-len(format.Tag)))
	emptyFile := buffer.Bytes()

	emptyFile[0x10] = 0x80
	emptyFile[0x11] = 0x40
	emptyFile[0x12] = 0x03
	emptyFile[0x18] = 0x80
	emptyFile[0x19] = 0x02
	emptyFile[0x1A] = 0xE0
	emptyFile[0x1B] = 0x01
	emptyFile[0x26] = 0x22
	emptyFile[0x27] = 0x56

	source := bytes.NewReader(emptyFile)
	container, _ := Read(source)

	c.Check(container.VideoWidth(), check.Equals, uint16(640))
	c.Check(container.VideoHeight(), check.Equals, uint16(480))
	c.Check(container.AudioSampleRate(), check.Equals, uint16(22050))
	c.Check(container.EntryCount(), check.Equals, 0)
}

func (suite *ReadSuite) TestReadReturnsContainerWithDataEntriesExceptTerminator(c *check.C) {
	testData := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	buffer := bytes.NewBufferString(format.Tag)
	buffer.Write(make([]byte, 0x100+0x300-len(format.Tag)))
	buffer.Write(make([]byte, 0xC00))
	buffer.Write(testData)
	raw := buffer.Bytes()

	raw[0x04] = 2
	// size of index table
	raw[0x08] = 0x00
	raw[0x09] = 0x0C

	// index entry 0
	raw[0x0400+3] = 0x02
	raw[0x0400+4] = 0x00
	raw[0x0400+5] = 0x10
	// index entry 1
	raw[0x0408+3] = 0x00
	raw[0x0408+4] = byte(len(testData))
	raw[0x0408+5] = 0x10

	source := bytes.NewReader(raw)
	container, _ := Read(source)

	c.Check(container.EntryCount(), check.Equals, 1)
	c.Check(container.Entry(0).Type(), check.Equals, Audio)
	c.Check(container.Entry(0).Data(), check.DeepEquals, testData)
}
