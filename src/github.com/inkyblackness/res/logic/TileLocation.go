package logic

// TileLocation specifies the location of a tile in the map. It contains the
// X and Y values of the tile.
type TileLocation struct {
	x, y uint16
}

// AtTile returns a tile location with the specified coordinates.
func AtTile(x, y uint16) TileLocation {
	return TileLocation{x: x, y: y}
}

// XY returns the X and Y values of the location.
func (loc TileLocation) XY() (x, y uint16) {
	return loc.x, loc.y
}
