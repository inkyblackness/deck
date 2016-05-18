package data

import (
	"fmt"
)

// TileCoordinate describes the position in one dimension on the tile map. It contains the tile number and an
// offset within the given tile.
type TileCoordinate uint16

func (coord TileCoordinate) String() (result string) {
	result += fmt.Sprintf("%d:%d", coord.Tile(), coord.Offset())

	return
}

// Tile returns the tile number for this coordinate.
func (coord TileCoordinate) Tile() byte {
	return byte(uint16(coord) >> 8)
}

// Offset returns the offset within the tile for this coordinate.
func (coord TileCoordinate) Offset() byte {
	return byte(coord & 0x00FF)
}
