package core

import (
	"bytes"
	"encoding/binary"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/image"
	"github.com/inkyblackness/res/text"
	"github.com/inkyblackness/res/textprop"
	"github.com/inkyblackness/shocked-core/io"
	model "github.com/inkyblackness/shocked-model"
)

// Textures is the adapter of texture information.
type Textures struct {
	cybstrng   [model.LanguageCount]*io.DynamicChunkStore
	images     *io.DynamicChunkStore
	cp         text.Codepage
	properties textprop.Store
}

// NewTextures returns a new Textures instance, if possible.
func NewTextures(library io.StoreLibrary) (textures *Textures, err error) {
	var cybstrng [model.LanguageCount]*io.DynamicChunkStore
	var images *io.DynamicChunkStore
	var properties textprop.Store

	for i := 0; i < model.LanguageCount && err == nil; i++ {
		cybstrng[i], err = library.ChunkStore(localized[i].cybstrng)
	}
	if err == nil {
		images, err = library.ChunkStore("texture.res")
	}
	if err == nil {
		properties, err = library.TextpropStore("textprop.dat")
	}

	if err == nil {
		textures = &Textures{cybstrng: cybstrng, images: images, cp: text.DefaultCodepage(), properties: properties}
	}

	return
}

// TextureCount returns the number of available textures.
func (textures *Textures) TextureCount() int {
	return 273
}

// Image returns the bitmap of identified & sized texture.
func (textures *Textures) Image(index int, size model.TextureSize) (bmp image.Bitmap) {
	var resID res.ResourceID
	blockIndex := uint16(0)

	if size == model.TextureLarge {
		resID = res.ResourceID(0x03E8 + index)
	} else if size == model.TextureMedium {
		resID = res.ResourceID(0x02C3 + index)
	} else if size == model.TextureSmall {
		resID = res.ResourceID(0x004D)
		blockIndex = uint16(index)
	} else if size == model.TextureIcon {
		resID = res.ResourceID(0x004C)
		blockIndex = uint16(index)
	}
	holder := textures.images.Get(resID)
	var blockData []byte
	if holder != nil {
		blockData = holder.BlockData(blockIndex)
	}
	if len(blockData) > 0 {
		bmp, _ = image.Read(bytes.NewReader(blockData))
	}

	return
}

// SetImage requests to set the bitmap of identified & sized texture..
func (textures *Textures) SetImage(index int, size model.TextureSize, imgBitmap image.Bitmap) {
	writer := bytes.NewBuffer(nil)
	image.Write(writer, imgBitmap, image.UncompressedBitmap, false, 0)
	blockData := writer.Bytes()

	if size == model.TextureLarge {
		textures.images.Get(res.ResourceID(0x03E8+index)).SetBlockData(0, blockData)
	} else if size == model.TextureMedium {
		textures.images.Get(res.ResourceID(0x02C3+index)).SetBlockData(0, blockData)
	} else if size == model.TextureSmall {
		textures.images.Get(res.ResourceID(0x004D)).SetBlockData(uint16(index), blockData)
	} else if size == model.TextureIcon {
		textures.images.Get(res.ResourceID(0x004C)).SetBlockData(uint16(index), blockData)
	}
}

func (textures *Textures) rawProperties(index int) (entry textprop.Entry) {
	rawProperties := textures.properties.Get(uint32(index))
	reader := bytes.NewReader(rawProperties)

	binary.Read(reader, binary.LittleEndian, &entry)

	return
}

func (textures *Textures) setRawProperties(index int, entry textprop.Entry) {
	writer := bytes.NewBuffer(nil)

	binary.Write(writer, binary.LittleEndian, &entry)
	textures.properties.Put(uint32(index), writer.Bytes())
}

// Properties returns the texture properties of the identified texture.
func (textures *Textures) Properties(index int) model.TextureProperties {
	prop := model.TextureProperties{}
	rawProperties := textures.rawProperties(index)

	for i := 0; i < model.LanguageCount; i++ {
		names := textures.cybstrng[i].Get(res.ResourceID(0x086A))
		cantBeUseds := textures.cybstrng[i].Get(res.ResourceID(0x086B))

		prop.Name[i] = textures.decodeString(names.BlockData(uint16(index)))
		prop.CantBeUsed[i] = textures.decodeString(cantBeUseds.BlockData(uint16(index)))
	}
	prop.Climbable = boolAsPointer(rawProperties.IsClimbable())
	prop.TransparencyControl = intAsPointer(int(rawProperties.TransparencyControl))
	prop.AnimationGroup = intAsPointer(int(rawProperties.AnimationGroup))
	prop.AnimationIndex = intAsPointer(int(rawProperties.AnimationIndex))

	return prop
}

// SetProperties requests to update the properties of identified texture.
func (textures *Textures) SetProperties(index int, prop model.TextureProperties) {
	rawProperties := textures.rawProperties(index)

	for i := 0; i < model.LanguageCount; i++ {
		if prop.Name[i] != nil {
			names := textures.cybstrng[i].Get(res.ResourceID(0x086A))
			names.SetBlockData(uint16(index), textures.encodeString(prop.Name[i]))
		}
		if prop.CantBeUsed[i] != nil {
			cantBeUseds := textures.cybstrng[i].Get(res.ResourceID(0x086B))
			cantBeUseds.SetBlockData(uint16(index), textures.encodeString(prop.CantBeUsed[i]))
		}
	}
	if prop.Climbable != nil {
		rawProperties.Climbable = 0
		if *prop.Climbable {
			rawProperties.Climbable = 1
		}
	}
	if prop.TransparencyControl != nil {
		rawProperties.TransparencyControl = byte(*prop.TransparencyControl)
	}
	if prop.AnimationGroup != nil {
		rawProperties.AnimationGroup = byte(*prop.AnimationGroup)
	}
	if prop.AnimationIndex != nil {
		rawProperties.AnimationIndex = byte(*prop.AnimationIndex)
	}
	textures.setRawProperties(index, rawProperties)
}

func (textures *Textures) decodeString(data []byte) *string {
	value := textures.cp.Decode(data)

	return &value
}

func (textures *Textures) encodeString(value *string) []byte {
	data := textures.cp.Encode(*value)

	return data
}
