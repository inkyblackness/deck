package model

// LevelObjectProperties describe the mutable properties of a level object
type LevelObjectProperties struct {
	Subclass *int `json:"subclass"`
	Type     *int `json:"type"`

	TileX *int `json:"tileX"`
	FineX *int `json:"fineX"`
	TileY *int `json:"tileY"`
	FineY *int `json:"fineY"`
	Z     *int `json:"z"`

	ClassData []byte `json:"classData"`

	RotationX *int
	RotationY *int
	RotationZ *int
	Hitpoints *int
}
