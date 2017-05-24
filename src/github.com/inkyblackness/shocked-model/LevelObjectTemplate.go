package model

// LevelObjectTemplate contains the minimum information for a new level object.
type LevelObjectTemplate struct {
	Class    int
	Subclass int
	Type     int

	TileX int
	FineX int
	TileY int
	FineY int
	Z     int

	Hitpoints int
}
