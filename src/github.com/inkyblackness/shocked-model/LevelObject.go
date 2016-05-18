package model

type LevelObject struct {
	Identifiable

	Class    int `json:"class"`
	Subclass int `json:"subclass"`
	Type     int `json:"type"`

	BaseProperties LevelObjectBaseProperties `json:"properties"`

	Hacking LevelObjectHacking

	Links []Link `json:"links"`
}

type LevelObjectHacking struct {
	Unknown0013 []int
	Unknown0015 []int
	Unknown0017 []int

	ClassData []int
}
