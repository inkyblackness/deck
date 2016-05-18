package font

// Font describes a bitmap font. Monochrome fonts are stored as packed bits, color fonts use bytes as palette indices.
type Font interface {
	// IsMonochrome returns true if the glyphs are packed as bits within the bitmap.
	IsMonochrome() bool

	// BitmapWidth returns the width of the bitmap in bytes.
	BitmapWidth() int
	// BitmapHeight returns the height of the bitmap in bytes.
	BitmapHeight() int
	// Bitmap returns the actual bitmap of the font.
	Bitmap() []byte

	// FirstCharacter returns the index of the first available character of the font (inclusive).
	FirstCharacter() int
	// LastCharacter returns the index of the last available character of the font (inclusive).
	LastCharacter() int

	// GlyphXOffset returns the horizontal offset to the glyph for given character index.
	// The width of a glyph is equal to the distance to the next one. This function returns a valid
	// result for the next index after the last character.
	GlyphXOffset(index int) int
}
