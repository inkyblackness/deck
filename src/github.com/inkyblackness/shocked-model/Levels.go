package model

// Levels is the collection of currently available levels.
type Levels struct {
	Referable

	// List refers to the available levels.
	List []Level `json:"list"`
}
