package font

import (
	"bytes"
	"encoding/binary"
)

// Save encodes a bitmap font into a stream of bytes.
func Save(font Font) []byte {
	writer := bytes.NewBuffer(nil)
	var header Header
	offsets := font.LastCharacter() - font.FirstCharacter() + 2

	if font.IsMonochrome() {
		header.Type = Monochrome
	} else {
		header.Type = Color
	}
	header.FirstCharacter = uint16(font.FirstCharacter())
	header.LastCharacter = uint16(font.LastCharacter())
	header.XOffsetStart = uint32(HeaderSize)
	header.BitmapStart = header.XOffsetStart + uint32(offsets*2)
	header.Width = uint16(font.BitmapWidth())
	header.Height = uint16(font.BitmapHeight())

	binary.Write(writer, binary.LittleEndian, &header)
	for i := 0; i < offsets; i++ {
		binary.Write(writer, binary.LittleEndian, uint16(font.GlyphXOffset(i)))
	}

	writer.Write(font.Bitmap())

	return writer.Bytes()
}
