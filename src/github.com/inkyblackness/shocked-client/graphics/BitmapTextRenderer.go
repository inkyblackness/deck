package graphics

import (
	"github.com/inkyblackness/res/text"
	"github.com/inkyblackness/shocked-model"
)

type bitmapTextRenderer struct {
	cp                 text.Codepage
	font               model.Font
	bitmap             Bitmap
	lastCharacterIndex int
}

// NewBitmapTextRenderer returns a new text renderer for the given bitmap font.
func NewBitmapTextRenderer(font model.Font) TextRenderer {
	return &bitmapTextRenderer{
		cp:                 text.DefaultCodepage(),
		font:               font,
		bitmap:             BitmapFromRaw(font.Bitmap),
		lastCharacterIndex: font.FirstCharacter + len(font.GlyphXOffsets) - 1}
}

func (renderer *bitmapTextRenderer) Render(text string) Bitmap {
	var bmp Bitmap
	indexLines := renderer.mapCharactersToIndex(text)

	for _, line := range indexLines {
		lineWidth := 2
		for _, characterIndex := range line {
			charWidth := renderer.font.GlyphXOffsets[characterIndex+1] - renderer.font.GlyphXOffsets[characterIndex]
			lineWidth += charWidth
		}
		if bmp.Width < lineWidth {
			bmp.Width = lineWidth
		}
	}
	bmp.Height = renderer.bitmap.Height*len(indexLines) + 1 + len(indexLines)
	bmp.Pixels = make([]byte, bmp.Width*bmp.Height)
	for lineIndex, line := range indexLines {
		outStartY := 1 + lineIndex + renderer.bitmap.Height*lineIndex
		outStartX := 1
		for _, characterIndex := range line {
			inStartX := renderer.font.GlyphXOffsets[characterIndex]
			inEndX := renderer.font.GlyphXOffsets[characterIndex+1]
			charWidth := inEndX - inStartX
			for y := 0; y < renderer.bitmap.Height; y++ {
				inStartY := renderer.bitmap.Width * y
				copy(bmp.Pixels[bmp.Width*(outStartY+y)+outStartX:], renderer.bitmap.Pixels[inStartY+inStartX:inStartY+inEndX])
			}
			outStartX += charWidth
		}
	}
	if renderer.font.Monochrome {
		renderer.outline(bmp)
	}

	return bmp
}

func (renderer *bitmapTextRenderer) mapCharactersToIndex(text string) [][]int {
	lines := [][]int{}
	curLine := []int{}

	for _, character := range text {
		if character == '\n' {
			lines = append(lines, curLine)
			curLine = []int{}
		} else {
			cpIndex := int(renderer.cp.Encode(string(character))[0])
			if (cpIndex >= renderer.font.FirstCharacter) && (cpIndex < renderer.lastCharacterIndex) {
				curLine = append(curLine, cpIndex-renderer.font.FirstCharacter)
			}
		}
	}
	lines = append(lines, curLine)

	return lines
}

func (renderer *bitmapTextRenderer) outline(bmp Bitmap) {
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
