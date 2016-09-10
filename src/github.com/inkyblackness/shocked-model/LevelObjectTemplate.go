package model

// LevelObjectTemplate contains the minimum information for a new level object.
type LevelObjectTemplate struct {
	Class    int `json:"class"`
	Subclass int `json:"subclass"`
	Type     int `json:"type"`

	TileX int `json:"tileX"`
	FineX int `json:"fineX"`
	TileY int `json:"tileY"`
	FineY int `json:"fineY"`
	Z     int `json:"z"`
}
