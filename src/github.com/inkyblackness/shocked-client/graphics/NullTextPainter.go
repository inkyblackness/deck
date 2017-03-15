package graphics

type nullTextPainter struct {
}

// NewNullTextPainter returns a TextPainter instance that creates bitmaps
// with 2x2 pixel, which have the value 0x00.
func NewNullTextPainter() TextPainter {
	return &nullTextPainter{}
}

func (painter *nullTextPainter) Paint(text string) TextBitmap {
	return TextBitmap{
		Bitmap: Bitmap{
			Width:  2,
			Height: 2,
			Pixels: []byte{0x00, 0x00, 0x00, 0x00}},
		lineHeight: 2,
		offsets:    [][]int{[]int{0}}}
}
