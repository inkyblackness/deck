package movi

import (
	"bytes"

	check "gopkg.in/check.v1"
)

type WriteSuite struct {
}

var _ = check.Suite(&WriteSuite{})

func (suite *WriteSuite) TestWriteOfEmptyContainerCreatesMinimumSizeData(c *check.C) {
	builder := NewContainerBuilder()
	container := builder.Build()
	buffer := bytes.NewBuffer(nil)

	Write(buffer, container)

	c.Check(len(buffer.Bytes()), check.Equals, 0x0800)
}

func (suite *WriteSuite) TestWriteCanSaveEmptyContainer(c *check.C) {
	builder := NewContainerBuilder()
	container := builder.Build()
	buffer := bytes.NewBuffer(nil)

	Write(buffer, container)

	result, err := Read(bytes.NewReader(buffer.Bytes()))

	c.Assert(err, check.IsNil)
	c.Assert(result, check.NotNil)
	c.Check(result.EntryCount(), check.Equals, 0)
}

func (suite *WriteSuite) TestWriteSavesEntries(c *check.C) {
	dataBytes := []byte{0x01, 0x02, 0x03}
	builder := NewContainerBuilder()
	builder.AudioSampleRate(22050.0)
	builder.AddEntry(NewMemoryEntry(0.0, Audio, dataBytes))
	container := builder.Build()
	buffer := bytes.NewBuffer(nil)

	Write(buffer, container)

	result, err := Read(bytes.NewReader(buffer.Bytes()))

	c.Assert(err, check.IsNil)
	c.Assert(result, check.NotNil)
	c.Assert(result.EntryCount(), check.Equals, 1)
	c.Check(result.Entry(0).Data(), check.DeepEquals, dataBytes)
}

func (suite *WriteSuite) TestIndexTableSizeFor_ExistingSizes(c *check.C) {
	// These sample values are always the minimum and maximum amount of index entries
	// found for a given index size.
	c.Check(indexTableSizeFor(3), check.Equals, 0x0400)
	c.Check(indexTableSizeFor(127), check.Equals, 0x0400)

	c.Check(indexTableSizeFor(130), check.Equals, 0x0C00)
	c.Check(indexTableSizeFor(218), check.Equals, 0x0C00)

	c.Check(indexTableSizeFor(738), check.Equals, 0x1C00)
	c.Check(indexTableSizeFor(755), check.Equals, 0x1C00)

	c.Check(indexTableSizeFor(1475), check.Equals, 0x3400)
	c.Check(indexTableSizeFor(1523), check.Equals, 0x3400)
}
