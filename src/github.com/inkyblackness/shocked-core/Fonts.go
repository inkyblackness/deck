package core

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	resFont "github.com/inkyblackness/res/font"
	"github.com/inkyblackness/shocked-core/io"
	model "github.com/inkyblackness/shocked-model"
)

// Fonts represents the game fonts accessor
type Fonts struct {
	gamescr *io.DynamicChunkStore
}

// NewFonts returns a new instance of Fonts.
func NewFonts(library io.StoreLibrary) (fonts *Fonts, err error) {
	var gamescr *io.DynamicChunkStore

	gamescr, err = library.ChunkStore("gamescr.res")

	if err == nil {
		fonts = &Fonts{gamescr: gamescr}
	}

	return
}

// Font returns the font data for the identified font.
func (fonts *Fonts) Font(id res.ResourceID) (font *model.Font, err error) {
	fontChunk := fonts.gamescr.Get(id)
	if fontChunk.ContentType() == chunk.Font {
		fontBlockData := fontChunk.BlockData(0)
		var fontData resFont.Font
		fontData, err = resFont.Load(bytes.NewReader(fontBlockData))

		if err == nil {
			isMonochrome := fontData.IsMonochrome()
			if isMonochrome {
				fontData = resFont.EnsureColor(fontData, 1)
			}

			font = &model.Font{
				Monochrome: isMonochrome,
				Bitmap: model.RawBitmap{
					Width:  fontData.BitmapWidth(),
					Height: fontData.BitmapHeight(),
					Pixels: base64.StdEncoding.EncodeToString(fontData.Bitmap())},
				FirstCharacter: fontData.FirstCharacter(),
				GlyphXOffsets:  make([]int, fontData.LastCharacter()-fontData.FirstCharacter())}

			for charIndex := 0; charIndex < len(font.GlyphXOffsets); charIndex++ {
				font.GlyphXOffsets[charIndex] = fontData.GlyphXOffset(charIndex)
			}
		} else {
			err = fmt.Errorf("Failed to load font ID %v", id)
		}
	} else {
		err = fmt.Errorf("ID %v is not a font", id)
	}

	return
}
