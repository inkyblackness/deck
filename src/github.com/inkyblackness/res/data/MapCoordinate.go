package data

import (
	"fmt"
)

// MapCoordinate describes the position in one dimension on the tile map. It contains the tile number and an
// offset within the given tile.
type MapCoordinate uint16

// MapCoordinateOf returns a map coordinate with the given values.
func MapCoordinateOf(tile byte, offset byte) MapCoordinate {
	return MapCoordinate(uint16(tile)<<8 + uint16(offset))
}

func (coord MapCoordinate) String() (result string) {
	result += fmt.Sprintf("%d:%d", coord.Tile(), coord.Offset())

	return
}

// Tile returns the tile number for this coordinate.
func (coord MapCoordinate) Tile() byte {
	return byte(uint16(coord) >> 8)
}

// Offset returns the offset within the tile for this coordinate.
func (coord MapCoordinate) Offset() byte {
	return byte(coord & 0x00FF)
}
