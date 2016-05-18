package font

type simpleFont struct {
	monochrome bool

	bitmapWidth  int
	bitmapHeight int
	bitmap       []byte

	firstCharacter int
	lastCharacter  int
	xOffsets       []uint16
}

func newSimpleFont(monochrome bool, width, height int, firstCharacter, lastCharacter int) *simpleFont {
	return &simpleFont{
		monochrome:     monochrome,
		bitmapWidth:    width,
		bitmapHeight:   height,
		bitmap:         make([]byte, width*height),
		firstCharacter: firstCharacter,
		lastCharacter:  lastCharacter,
		xOffsets:       make([]uint16, lastCharacter-firstCharacter+2)}
}

func (font *simpleFont) IsMonochrome() bool {
	return font.monochrome
}

func (font *simpleFont) BitmapWidth() int {
	return font.bitmapWidth
}

func (font *simpleFont) BitmapHeight() int {
	return font.bitmapHeight
}

func (font *simpleFont) Bitmap() []byte {
	return font.bitmap
}

func (font *simpleFont) FirstCharacter() int {
	return font.firstCharacter
}

func (font *simpleFont) LastCharacter() int {
	return font.lastCharacter
}

func (font *simpleFont) GlyphXOffset(index int) int {
	return int(font.xOffsets[index])
}
