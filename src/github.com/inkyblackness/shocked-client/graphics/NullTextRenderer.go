package graphics

type nullTextRenderer struct {
}

// NewNullTextRenderer returns a TextRenderer instance that creates bitmaps
// with one pixel, which has the value 0x00.
func NewNullTextRenderer() TextRenderer {
	return &nullTextRenderer{}
}

func (renderer *nullTextRenderer) Render(text string) Bitmap {
	return Bitmap{
		Width:  2,
		Height: 2,
		Pixels: []byte{0x00, 0x00, 0x00, 0x00}}
}
