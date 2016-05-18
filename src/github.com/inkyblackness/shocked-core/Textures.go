package core

import (
	"bytes"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/image"
	"github.com/inkyblackness/res/text"
	"github.com/inkyblackness/shocked-core/io"
	model "github.com/inkyblackness/shocked-model"
)

type Textures struct {
	cybstrng [model.LanguageCount]chunk.Store
	images   chunk.Store
	cp       text.Codepage
}

func NewTextures(library io.StoreLibrary) (textures *Textures, err error) {
	var cybstrng [model.LanguageCount]chunk.Store
	var images chunk.Store

	for i := 0; i < model.LanguageCount && err == nil; i++ {
		cybstrng[i], err = library.ChunkStore(localized[i].cybstrng)
	}
	if err == nil {
		images, err = library.ChunkStore("texture.res")
	}

	if err == nil {
		textures = &Textures{cybstrng: cybstrng, images: images, cp: text.DefaultCodepage()}
	}

	return
}

func (textures *Textures) TextureCount() int {
	return 273
}

func (textures *Textures) Image(index int, size model.TextureSize) (bmp image.Bitmap) {
	var blockData []byte

	if size == model.TextureLarge {
		blockData = textures.images.Get(res.ResourceID(0x03E8 + index)).BlockData(0)
	} else if size == model.TextureMedium {
		blockData = textures.images.Get(res.ResourceID(0x02C3 + index)).BlockData(0)
	} else if size == model.TextureSmall {
		blockData = textures.images.Get(res.ResourceID(0x004D)).BlockData(uint16(index))
	} else if size == model.TextureIcon {
		blockData = textures.images.Get(res.ResourceID(0x004C)).BlockData(uint16(index))
	}
	bmp, _ = image.Read(bytes.NewReader(blockData))

	return
}

func (textures *Textures) Properties(index int) model.TextureProperties {
	prop := model.TextureProperties{}

	for i := 0; i < model.LanguageCount; i++ {
		names := textures.cybstrng[i].Get(res.ResourceID(0x086A))
		cantBeUseds := textures.cybstrng[i].Get(res.ResourceID(0x086B))

		prop.Name[i] = textures.DecodeString(names.BlockData(uint16(index)))
		prop.CantBeUsed[i] = textures.DecodeString(cantBeUseds.BlockData(uint16(index)))
	}

	return prop
}

func (textures *Textures) SetProperties(index int, prop model.TextureProperties) {
	for i := 0; i < model.LanguageCount; i++ {
		if prop.Name[i] != nil {
			names := textures.cybstrng[i].Get(res.ResourceID(0x086A))
			names.SetBlockData(uint16(index), textures.EncodeString(prop.Name[i]))
		}
		if prop.CantBeUsed[i] != nil {
			cantBeUseds := textures.cybstrng[i].Get(res.ResourceID(0x086B))
			cantBeUseds.SetBlockData(uint16(index), textures.EncodeString(prop.CantBeUsed[i]))
		}
	}
}

func (textures *Textures) DecodeString(data []byte) *string {
	value := textures.cp.Decode(data)

	return &value
}

func (textures *Textures) EncodeString(value *string) []byte {
	data := textures.cp.Encode(*value)

	return data
}
