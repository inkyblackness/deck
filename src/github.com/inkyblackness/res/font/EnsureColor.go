package font

// EnsureColor returns a color font based on the provided one. If the provided font is already a colored font,
// this one is returned. For a monochrome font, a new one will be created with the given color value set for the glyphs.
func EnsureColor(font Font, color byte) Font {
	var result Font

	if font.IsMonochrome() {
		simple := newSimpleFont(false, font.BitmapWidth()*8, font.BitmapHeight(),
			font.FirstCharacter(), font.LastCharacter())

		for i := 0; i < len(simple.xOffsets); i++ {
			simple.xOffsets[i] = uint16(font.GlyphXOffset(i))
		}
		for inIndex, mask := range font.Bitmap() {
			outIndex := inIndex * 8
			for offset := uint(0); offset < 8; offset++ {
				if (mask & (0x80 >> offset)) != 0 {
					simple.bitmap[outIndex+int(offset)] = color
				}
			}
		}
		result = simple
	} else {
		result = font
	}

	return result
}
