package model

// LevelObject describes one object of the level.
type LevelObject struct {
	Identifiable

	Class    int `json:"class"`
	Subclass int `json:"subclass"`
	Type     int `json:"type"`

	BaseProperties LevelObjectBaseProperties `json:"properties"`

	ClassData []int `json:"classData"`

	Links []Link `json:"links"`
}
