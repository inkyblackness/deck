package font

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// Load decodes a bitmap font from given source. It returns the font on success or an error otherwise.
func Load(source io.ReadSeeker) (font Font, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	if source == nil {
		panic(fmt.Errorf("source is nil"))
	}

	var header Header
	binary.Read(source, binary.LittleEndian, &header)

	simpleFont := newSimpleFont(header.Type == Monochrome,
		int(header.Width), int(header.Height),
		int(header.FirstCharacter), int(header.LastCharacter))

	source.Seek(int64(header.XOffsetStart), os.SEEK_SET)
	binary.Read(source, binary.LittleEndian, simpleFont.xOffsets)

	source.Seek(int64(header.BitmapStart), os.SEEK_SET)
	source.Read(simpleFont.bitmap)

	font = simpleFont

	return
}
