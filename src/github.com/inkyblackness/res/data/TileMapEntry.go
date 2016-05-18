package data

import (
	"fmt"
)

// TileMapEntry describes one tile in the map.
type TileMapEntry struct {
	Type             TileType
	Floor            byte
	Ceiling          byte
	SlopeHeight      byte
	FirstObjectIndex uint16
	Textures         TileTextureInfo
	Flags            TileFlag
	UnknownState     [4]byte
}

// DefaultTileMapEntry returns a TileMapEntry instance with default values set.
func DefaultTileMapEntry() *TileMapEntry {
	entry := &TileMapEntry{
		Type:         Solid,
		UnknownState: [4]byte{0xFF, 0x00, 0x00, 0x00}}

	return entry
}

func (entry *TileMapEntry) String() (result string) {
	result += fmt.Sprintf("Type: %v\n", entry.Type)
	result += fmt.Sprintf("First Object Index: %d\n", entry.FirstObjectIndex)
	result += fmt.Sprintf("Flags: 0x%08X (%v)\n", uint32(entry.Flags), entry.Flags)

	return
}
