package model

// LevelObject describes one object of the level.
type LevelObject struct {
	Identifiable

	Class      int                   `json:"class"`
	Properties LevelObjectProperties `json:"properties"`

	Links []Link `json:"links"`
}
