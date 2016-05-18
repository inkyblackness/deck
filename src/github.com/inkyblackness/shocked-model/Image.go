package model

// Image describes a graphical facette of a larger entity. The image has some
// properties, together with format(s) for visual representation.
type Image struct {
	Referable

	// Properties consist of extra meta-information about an image
	Properties ImageProperties `json:"properties"`

	// Formats is a list of links with available formats
	Formats []Link `json:"formats"`
}
