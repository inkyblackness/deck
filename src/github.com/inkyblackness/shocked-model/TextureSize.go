package model

// TextureSize is an enumeration of allowed texture sizes.
type TextureSize string

const (
	TextureIcon   TextureSize = "icon"
	TextureSmall              = "small"
	TextureMedium             = "medium"
	TextureLarge              = "large"
)

// TextureSizes returns an array of all supported texture sizes
func TextureSizes() []TextureSize {
	return []TextureSize{TextureIcon, TextureSmall, TextureMedium, TextureLarge}
}
