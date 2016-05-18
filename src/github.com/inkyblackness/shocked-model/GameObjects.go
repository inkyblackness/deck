package model

// GameObjects is the collection of all game objects.
type GameObjects struct {
	Referable

	// Items refers to the available objects.
	Items []Identifiable `json:"items"`
}
