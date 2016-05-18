package model

// Color describes a single color value. The channels have a range of 0..255 .
type Color struct {
	Red   int `json:"r"`
	Green int `json:"g"`
	Blue  int `json:"b"`
}

// Palette is an identifiable list of colors.
type Palette struct {
	Identifiable

	// Colors contains the color values of the palette.
	Colors [256]Color `json:"colors"`
}
