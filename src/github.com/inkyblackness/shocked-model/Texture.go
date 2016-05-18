package model

// Texture is an identifiable combination of properties and the list of images
// for that texture.
type Texture struct {
	Identifiable

	// Properties contains the behavioural settings of a texture
	Properties TextureProperties `json:"properties"`

	// Images is a list of links to the associated images.
	Images []Link `json:"images"`
}
