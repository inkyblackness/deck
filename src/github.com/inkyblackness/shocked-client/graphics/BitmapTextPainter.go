package graphics

import (
	"strings"

	"github.com/inkyblackness/res/text"
	"github.com/inkyblackness/shocked-model"
)

type bitmapTextPainter struct {
	cp                 text.Codepage
	font               model.Font
	bitmap             Bitmap
	lastCharacterIndex int
}

// NewBitmapTextPainter returns a new text painter for the given bitmap font.
func NewBitmapTextPainter(font model.Font) TextPainter {
	return &bitmapTextPainter{
		cp:                 text.DefaultCodepage(),
		font:               font,
		bitmap:             BitmapFromRaw(font.Bitmap),
		lastCharacterIndex: font.FirstCharacter + len(font.GlyphXOffsets) - 1}
}

func (painter *bitmapTextPainter) Paint(text string, widthLimit int) TextBitmap {
	var bmp TextBitmap
	indexLines := painter.mapCharactersToIndex(text, widthLimit)

	bmp.lineHeight = painter.bitmap.Height + 1
	for _, line := range indexLines {
		lineWidth := 2
		lineOffsets := []int{0}
		for characterOffset, characterIndex := range line {
			charWidth := painter.font.GlyphXOffsets[characterIndex+1] - painter.font.GlyphXOffsets[characterIndex]
			lineWidth += charWidth
			lineOffsets = append(lineOffsets, lineOffsets[characterOffset]+charWidth)
		}
		bmp.offsets = append(bmp.offsets, lineOffsets)
		if bmp.Width < lineWidth {
			bmp.Width = lineWidth
		}
	}
	bmp.Height = painter.bitmap.Height*len(indexLines) + 1 + len(indexLines)
	bmp.Pixels = make([]byte, bmp.Width*bmp.Height)
	for lineIndex, line := range indexLines {
		outStartY := 1 + lineIndex + painter.bitmap.Height*lineIndex
		outStartX := 1
		for _, characterIndex := range line {
			inStartX := painter.font.GlyphXOffsets[characterIndex]
			inEndX := painter.font.GlyphXOffsets[characterIndex+1]
			charWidth := inEndX - inStartX
			for y := 0; y < painter.bitmap.Height; y++ {
				inStartY := painter.bitmap.Width * y
				copy(bmp.Pixels[bmp.Width*(outStartY+y)+outStartX:], painter.bitmap.Pixels[inStartY+inStartX:inStartY+inEndX])
			}
			outStartX += charWidth
		}
	}
	if painter.font.Monochrome {
		painter.outline(bmp.Bitmap)
	}

	return bmp
}

type indexWord struct {
	totalSize int
	indices   []int
}

func (word indexWord) join(other indexWord) indexWord {
	joinedIndices := make([]int, len(word.indices)+len(other.indices))
	copy(joinedIndices, word.indices)
	copy(joinedIndices[len(word.indices):], other.indices)

	return indexWord{totalSize: word.totalSize + other.totalSize, indices: joinedIndices}
}

func (painter *bitmapTextPainter) mapCharactersToIndex(text string, widthLimit int) [][]int {
	result := [][]int{}
	toIndexWord := func(word string) (result indexWord) {
		encoded := painter.cp.Encode(word)
		for _, charIndexByte := range encoded {
			charIndex := int(charIndexByte)
			if (charIndex >= painter.font.FirstCharacter) && (charIndex < painter.lastCharacterIndex) {
				offsetIndex := charIndex - painter.font.FirstCharacter
				result.indices = append(result.indices, offsetIndex)
				result.totalSize += painter.font.GlyphXOffsets[offsetIndex+1] - painter.font.GlyphXOffsets[offsetIndex]
			}
		}
		return
	}
	textLines := strings.Split(text, "\n")
	resultLine := toIndexWord("")
	newLine := func() {
		result = append(result, resultLine.indices)
		resultLine = toIndexWord("")
	}
	space := toIndexWord(" ")

	for _, textLine := range textLines {
		words := strings.Split(textLine, " ")
		for wordIndex, word := range words {
			iWord := toIndexWord(word)
			if wordIndex > 0 {
				resultLine = resultLine.join(space)
			}
			if (widthLimit > 0) && ((resultLine.totalSize + iWord.totalSize) > widthLimit) {
				newLine()
			}
			resultLine = resultLine.join(iWord)
		}
		newLine()
	}

	return result
}

func (painter *bitmapTextPainter) outline(bmp Bitmap) {
	perimeter := func(index, limit int) (values []int) {
		if index > 0 {
			values = append(values, -1)
		}
		values = append(values, 0)
		if index < (limit - 1) {
			values = append(values, 1)
		}
		return
	}

	for pixelOffset, pixelValue := range bmp.Pixels {
		if pixelValue == 0 {
			lines := perimeter(pixelOffset/bmp.Width, bmp.Height)
			columns := perimeter(pixelOffset%bmp.Width, bmp.Width)
			isNeighbour := false

			for _, lineOffset := range lines {
				for _, columnOffset := range columns {
					if !isNeighbour && (bmp.Pixels[pixelOffset+lineOffset*bmp.Width+columnOffset] == 1) {
						isNeighbour = true
					}
				}
			}
			if isNeighbour {
				bmp.Pixels[pixelOffset] = 2
			}
		}
	}
}
