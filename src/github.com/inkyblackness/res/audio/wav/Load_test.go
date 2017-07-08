package wav

import (
	"bytes"

	check "gopkg.in/check.v1"
)

type LoadSuite struct {
}

var _ = check.Suite(&LoadSuite{})

func (suite *LoadSuite) TestLoadReturnsErrorOnNil(c *check.C) {
	_, err := Load(nil)

	c.Check(err, check.NotNil)
}

func (suite *LoadSuite) TestLoadExtractsDataOfL8(c *check.C) {
	input := []byte{
		0x52, 0x49, 0x46, 0x46, // "RIFF"
		0x29, 0x00, 0x00, 0x00, // len(RIFF)
		0x57, 0x41, 0x56, 0x45, // "WAVE"
		0x66, 0x6d, 0x74, 0x20, // "fmt "
		0x10, 0x00, 0x00, 0x00, // len(fmt)
		0x01, 0x00, // fmt:type
		0x01, 0x00, // fmt:channels
		0x22, 0x56, 0x00, 0x00, // fmt:samples/sec
		0x22, 0x56, 0x00, 0x00, // fmt:avgBytes/sec
		0x01, 0x00, // fmt:blockAlign
		0x08, 0x00, // fmt:bits/sample
		0x64, 0x61, 0x74, 0x61, // "data"
		0x05, 0x00, 0x00, 0x00, // len(data)
		0x00, 0x40, 0x80, 0xC0, 0xFF} // data

	data, err := Load(bytes.NewReader(input))

	c.Assert(err, check.IsNil)
	c.Check(uint32(data.SampleRate()), check.Equals, uint32(22050))
	c.Check(data.Samples(0, data.SampleCount()), check.DeepEquals, []byte{0x00, 0x40, 0x80, 0xC0, 0xFF})
}

func (suite *LoadSuite) TestLoadExtractsDataOfL16(c *check.C) {
	input := []byte{
		0x52, 0x49, 0x46, 0x46, // "RIFF"
		0x2E, 0x00, 0x00, 0x00, // len(RIFF)
		0x57, 0x41, 0x56, 0x45, // "WAVE"
		0x66, 0x6d, 0x74, 0x20, // "fmt "
		0x10, 0x00, 0x00, 0x00, // len(fmt)
		0x01, 0x00, // fmt:type
		0x01, 0x00, // fmt:channels
		0x22, 0x56, 0x00, 0x00, // fmt:samples/sec
		0x44, 0xAC, 0x00, 0x00, // fmt:avgBytes/sec
		0x02, 0x00, // fmt:blockAlign
		0x10, 0x00, // fmt:bits/sample
		0x64, 0x61, 0x74, 0x61, // "data"
		0x0A, 0x00, 0x00, 0x00, // len(data)
		0xAA, 0x00, 0xAA, 0x40, 0xAA, 0x7F, 0xAA, 0xC0, 0xAA, 0xFF} // data

	data, err := Load(bytes.NewReader(input))

	c.Assert(err, check.IsNil)
	c.Check(uint32(data.SampleRate()), check.Equals, uint32(22050))
	c.Check(data.Samples(0, data.SampleCount()), check.DeepEquals, []byte{0x80, 0xC0, 0xFF, 0x40, 0x7F})
}
