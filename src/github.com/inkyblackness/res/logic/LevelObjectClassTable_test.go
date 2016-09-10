package logic

import (
	"github.com/inkyblackness/res/data"

	check "gopkg.in/check.v1"
)

type LevelObjectClassTableSuite struct {
}

var _ = check.Suite(&LevelObjectClassTableSuite{})

func (suite *LevelObjectClassTableSuite) TestEncodeSerializesAllData(c *check.C) {
	table := NewLevelObjectClassTable(data.LevelObjectPrefixSize+2, 2)

	entry0 := table.Entry(0)
	entry0.LevelObjectTableIndex = 0xABCD
	entry0.Previous = 0xEF44
	entry0.Next = 0x5566
	copy(entry0.Data(), []byte{0x01, 0x02})

	entry1 := table.Entry(1)
	copy(entry1.Data(), []byte{0x11, 0x12})

	serialized := table.Encode()

	c.Check(serialized, check.DeepEquals,
		[]byte{0xCD, 0xAB, 0x44, 0xEF, 0x66, 0x55, 0x1, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x11, 0x12})
}

func (suite *LevelObjectClassTableSuite) TestDecodeSerializesAllData(c *check.C) {
	entrySize := data.LevelObjectPrefixSize + 2
	table := NewLevelObjectClassTable(entrySize, 2)

	entry0 := table.Entry(0)
	entry0.LevelObjectTableIndex = 0xABCD
	entry0.Previous = 0xEF44
	entry0.Next = 0x5566
	copy(entry0.Data(), []byte{0x01, 0x02})

	entry1 := table.Entry(1)
	copy(entry1.Data(), []byte{0x11, 0x12})

	serialized := table.Encode()

	newTable := DecodeLevelObjectClassTable(serialized, entrySize)

	c.Check(newTable, check.DeepEquals, table)
}

func (suite *LevelObjectClassTableSuite) TestAsChainReturnsAChainViewOnTheTable(c *check.C) {
	entrySize := data.LevelObjectPrefixSize + 1
	table := NewLevelObjectClassTable(entrySize, 3)

	table.AsChain().Initialize(table.Count() - 1)

	checked := 1
	previous := table.Entry(0).Previous
	for previous != 0 && checked < 10 {
		checked++
		previous = table.Entry(data.LevelObjectChainIndex(previous)).Previous
	}

	c.Check(checked, check.Equals, table.Count())
}
