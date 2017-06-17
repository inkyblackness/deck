package model

// CalculatedWallHeights describes for each direction of a tile how walled off
// the given side is of a tile.
type CalculatedWallHeights struct {
	// North side of the tile (up)
	North float32
	// East side of the tile (right)
	East float32
	// West side of the tile (left)
	West float32
	// South side of the tile (down)
	South float32
}
