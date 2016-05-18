package model

// Workspace is the root object of the shocked model. There is only one workspace,
// which is why it has no dedicated ID.
type Workspace struct {
	Referable

	// Projects refers to the list of projects
	Projects Referable `json:"projects"`
}
