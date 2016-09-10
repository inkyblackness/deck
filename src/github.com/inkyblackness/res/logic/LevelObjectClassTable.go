package logic

import (
	"bytes"
	"encoding/binary"

	"github.com/inkyblackness/res/data"
)

// LevelObjectClassTable is a list of generic class entries.
type LevelObjectClassTable struct {
	entries []*data.LevelObjectClassEntry
}

// NewLevelObjectClassTable returns a new table of given amount of entries, each with
// the requested size. The size includes the LevelObjectPrefix size.
func NewLevelObjectClassTable(entrySize int, entryCount int) *LevelObjectClassTable {
	entries := make([]*data.LevelObjectClassEntry, entryCount)

	for index := 0; index < entryCount; index++ {
		entries[index] = data.NewLevelObjectClassEntry(entrySize - data.LevelObjectPrefixSize)
	}

	return &LevelObjectClassTable{entries: entries}
}

// DecodeLevelObjectClassTable reads a class table with entries of given size from the serialized
// byte array.
func DecodeLevelObjectClassTable(serialized []byte, entrySize int) *LevelObjectClassTable {
	entryCount := len(serialized) / entrySize
	table := NewLevelObjectClassTable(entrySize, entryCount)
	buf := bytes.NewReader(serialized)

	for _, entry := range table.entries {
		binary.Read(buf, binary.LittleEndian, &entry.LevelObjectPrefix)
		buf.Read(entry.Data())
	}

	return table
}

// Count returns the amount of entries in the table.
func (table *LevelObjectClassTable) Count() int {
	return len(table.entries)
}

// Entry returns the reference to a specified entry of the table.
func (table *LevelObjectClassTable) Entry(index data.LevelObjectChainIndex) *data.LevelObjectClassEntry {
	return table.entries[index]
}

// Encode serializes the table to an array of bytes.
func (table *LevelObjectClassTable) Encode() []byte {
	buf := bytes.NewBuffer(nil)

	for _, entry := range table.entries {
		binary.Write(buf, binary.LittleEndian, &entry.LevelObjectPrefix)
		buf.Write(entry.Data())
	}

	return buf.Bytes()
}

// AsChain returns a view on the table, interpreting it as a level object chain.
func (table *LevelObjectClassTable) AsChain() *LevelObjectChain {
	return NewLevelObjectChain(table.Entry(data.LevelObjectChainStartIndex),
		func(index data.LevelObjectChainIndex) LevelObjectChainLink {
			return table.Entry(index)
		})
}
