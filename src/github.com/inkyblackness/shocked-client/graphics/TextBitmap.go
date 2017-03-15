package graphics

import (
	"math"
)

// TextBitmap is a bitmap with text coordinates.
type TextBitmap struct {
	Bitmap

	lineHeight int
	offsets    [][]int
}

// LineHeight returns the height of one line, in pixel.
func (bmp TextBitmap) LineHeight() int {
	return bmp.lineHeight
}

// LineCount returns the number of lines in this bitmap.
func (bmp TextBitmap) LineCount() int {
	return len(bmp.offsets)
}

// LineLength returns the width of the given line, in pixel.
func (bmp TextBitmap) LineLength(line int) int {
	return bmp.CharOffset(line, math.MaxInt32)
}

// CharOffset returns the horizontal offset of given char in given line, in pixel.
// An unknown line results in an offset of zero, a char beyond the end results in the line's end.
func (bmp TextBitmap) CharOffset(line, char int) int {
	offset := 0

	if line < bmp.LineCount() {
		lineOffsets := bmp.offsets[line]
		chars := len(lineOffsets)
		if char < chars {
			offset = lineOffsets[char]
		} else {
			offset = lineOffsets[chars-1]
		}
	}

	return offset
}
