package model

// LevelObjectProperties describe the mutable properties of a level object
type LevelObjectProperties struct {
	Subclass *int
	Type     *int

	TileX *int
	FineX *int
	TileY *int
	FineY *int
	Z     *int

	ClassData []byte

	RotationX *int
	RotationY *int
	RotationZ *int
	Hitpoints *int

	ExtraData []byte
}
