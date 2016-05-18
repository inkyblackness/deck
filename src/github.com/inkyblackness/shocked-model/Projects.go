package model

// Projects is the collection of currently available projects.
type Projects struct {
	Referable

	// Items refers to the available projects.
	Items []Identifiable `json:"items"`
}
