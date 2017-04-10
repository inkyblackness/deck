package model

// GameObject describes one general game object
type GameObject struct {
	Identifiable

	Class    int `json:"class"`
	Subclass int `json:"subclass"`
	Type     int `json:"type"`

	Properties GameObjectProperties `json:"properties"`
}
