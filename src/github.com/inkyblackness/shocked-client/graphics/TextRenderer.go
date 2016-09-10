package graphics

// TextRenderer creates bitmaps for texts.
type TextRenderer interface {
	// Render creates a new bitmap based on a text.
	Render(text string) Bitmap
}
