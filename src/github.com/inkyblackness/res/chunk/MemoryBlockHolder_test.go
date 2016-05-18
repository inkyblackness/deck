package chunk

import (
	"github.com/inkyblackness/res"

	check "gopkg.in/check.v1"
)

type MemoryBlockHolderSuite struct {
}

var _ = check.Suite(&MemoryBlockHolderSuite{})

func (suite *MemoryBlockHolderSuite) TestChunkTypeReturnsProvidedValue(c *check.C) {
	holder := NewBlockHolder(BasicChunkType.WithCompression(), res.Palette, [][]byte{})

	c.Assert(holder.ChunkType(), check.Equals, BasicChunkType.WithCompression())
}

func (suite *MemoryBlockHolderSuite) TestContentTypeReturnsProvidedValue(c *check.C) {
	holder := NewBlockHolder(BasicChunkType.WithCompression(), res.Palette, [][]byte{})

	c.Assert(holder.ContentType(), check.Equals, res.Palette)
}

func (suite *MemoryBlockHolderSuite) TestBlockCountReturnsNumberOfBlocks(c *check.C) {
	holder := NewBlockHolder(BasicChunkType.WithCompression(), res.Palette, [][]byte{[]byte{0x01}, []byte{0x02}})

	c.Assert(holder.BlockCount(), check.Equals, uint16(2))
}

func (suite *MemoryBlockHolderSuite) TestBlockDataReturnsData(c *check.C) {
	holder := NewBlockHolder(BasicChunkType.WithCompression(), res.Palette, [][]byte{[]byte{0x01}, []byte{0x02}})

	c.Assert(holder.BlockData(1), check.DeepEquals, []byte{0x02})
}
