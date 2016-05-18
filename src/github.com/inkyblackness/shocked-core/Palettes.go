package core

import (
	"bytes"
	"image/color"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/image"

	"github.com/inkyblackness/shocked-core/io"
)

type Palettes struct {
	gamepal chunk.Store
}

func NewPalettes(library io.StoreLibrary) (palettes *Palettes, err error) {
	var gamepal chunk.Store

	gamepal, err = library.ChunkStore("gamepal.res")

	if err == nil {
		palettes = &Palettes{gamepal: gamepal}
	}

	return
}

func (palettes *Palettes) GamePalette() (color.Palette, error) {
	blockData := palettes.gamepal.Get(res.ResourceID(0x02BC)).BlockData(0)

	return image.LoadPalette(bytes.NewReader(blockData))
}
