package model

// Font describes an ingame bitmap font.
type Font struct {
	// Monochrome fonts have black/white bitmaps. They receive color later.
	Monochrome bool
	// Bitmap contains the raw bitmap of the font. Monochrome fonts have their pixel values set to 1 for visible pixels.
	Bitmap RawBitmap

	// FirstCharacter is the index of the first represented character of this font.
	FirstCharacter int
	// GlyphXOffsets is the horizontal offset for the character with given index. The width of a
	// character is the different to the next characters offset.
	GlyphXOffsets []int
}
