package model

// GameObject describes one general game object
type GameObject struct {
	Identifiable

	Properties GameObjectProperties `json:"properties"`
}
