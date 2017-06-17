package model

// GameObject describes one general game object
type GameObject struct {
	Class    int
	Subclass int
	Type     int

	Properties GameObjectProperties
}
