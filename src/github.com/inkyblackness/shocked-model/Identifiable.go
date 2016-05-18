package model

// Identifiable is something that can be referred to and be identified explicitly.
// Most of the things with an ID also can be addressed.
type Identifiable struct {
	Referable

	// ID is the unique identifier of the object
	ID string `json:"id"`
}
