package model

// Link describes one reference.
type Link struct {
	// Rel describes the relationship to the referred document.
	Rel string `json:"rel"`
	// Href holds the hyperlink to the referred document.
	Href string `json:"href"`
}
