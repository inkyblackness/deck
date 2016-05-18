package chunk

import (
	check "gopkg.in/check.v1"
)

type TypeIDSuite struct {
}

var _ = check.Suite(&TypeIDSuite{})

func (suite *TypeIDSuite) TestBasicChunkTypeIsNotCompressed(c *check.C) {
	c.Check(BasicChunkType.IsCompressed(), check.Equals, false)
}

func (suite *TypeIDSuite) TestBasicChunkTypeHasNoDirectory(c *check.C) {
	c.Check(BasicChunkType.HasDirectory(), check.Equals, false)
}

func (suite *TypeIDSuite) TestChunkTypeWithCompression(c *check.C) {
	c.Check(BasicChunkType.WithCompression().IsCompressed(), check.Equals, true)
}

func (suite *TypeIDSuite) TestChunkTypeWithDirectory(c *check.C) {
	c.Check(BasicChunkType.WithDirectory().HasDirectory(), check.Equals, true)
}

func (suite *TypeIDSuite) TestChunkTypeWithoutCompression(c *check.C) {
	c.Check(BasicChunkType.WithCompression().WithoutCompression().IsCompressed(), check.Equals, false)
}

func (suite *TypeIDSuite) TestChunkTypeWithoutDirectory(c *check.C) {
	c.Check(BasicChunkType.WithDirectory().WithoutDirectory().HasDirectory(), check.Equals, false)
}
