package graphics

// TextPainter creates bitmaps for texts.
type TextPainter interface {
	// Paint creates a new bitmap based on given text.
	Paint(text string, widthLimit int) TextBitmap
}
