package model

import (
	"fmt"
)

// TileCoordinate describes a location within a tile map.
type TileCoordinate struct {
	x, y int
}

// TileCoordinateOf returns a TileCoordinate with the given values.
func TileCoordinateOf(x, y int) TileCoordinate {
	return TileCoordinate{x: x, y: y}
}

// XY returns the x and y values.
func (coord TileCoordinate) XY() (x int, y int) {
	return coord.x, coord.y
}

// String implements the Stringer interface.
func (coord TileCoordinate) String() string {
	return fmt.Sprintf("%2d/%2d", coord.x, coord.y)
}
