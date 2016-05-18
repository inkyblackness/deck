package editor

const (
	// ZoomLevelMin specifies the minimum zoom level (farthest zoom).
	ZoomLevelMin = -2
	// ZoomLevelMax specifies the maximum zoom level (closest zoom).
	ZoomLevelMax = 4

	// TileBaseLength specifies the side length of a tile in world coordinates.
	TileBaseLength = float32(32)
	// TilesPerMapSide specifies how many tiles are per side of a map.
	// This value is implicitly identical to the actual data.
	TilesPerMapSide = 64
)
