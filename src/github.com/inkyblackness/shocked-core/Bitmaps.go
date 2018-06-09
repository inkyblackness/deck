package core

import (
	"bytes"
	"fmt"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/image"
	"github.com/inkyblackness/shocked-core/io"
	model "github.com/inkyblackness/shocked-model"
)

// Bitmaps is the adapter for general bitmaps.
type Bitmaps struct {
	mfdArt [model.LanguageCount]*io.DynamicChunkStore
}

// NewBitmaps returns a new Bitmaps instance, if possible.
func NewBitmaps(library io.StoreLibrary) (bitmaps *Bitmaps, err error) {
	var mfdArt [model.LanguageCount]*io.DynamicChunkStore

	for i := 0; i < model.LanguageCount && err == nil; i++ {
		mfdArt[i], err = library.ChunkStore(localized[i].mfdart)
	}

	if err == nil {
		bitmaps = &Bitmaps{mfdArt: mfdArt}
	}

	return
}

// Image returns the image data of identified bitmap.
func (bitmaps *Bitmaps) Image(key model.ResourceKey) (bmp image.Bitmap, err error) {
	var blockData []byte

	if (key.Type == model.ResourceTypeMfdDataImages) && key.HasValidLanguage() {
		holder := bitmaps.mfdArt[key.Language.ToIndex()].Get(res.ResourceID(key.Type))
		if key.Index < holder.BlockCount() {
			blockData = holder.BlockData(key.Index)
		}
	} else {
		err = fmt.Errorf("Unsupported resource key: %v", key)
	}

	if (err == nil) && (len(blockData) > 0) {
		bmp, err = image.Read(bytes.NewReader(blockData))
	} else {
		bmp = image.NullBitmap()
	}

	return
}

// SetImage requests to set the bitmap data of a resource.
func (bitmaps *Bitmaps) SetImage(key model.ResourceKey, bmp image.Bitmap) (resultKey model.ResourceKey, err error) {
	if (key.Type == model.ResourceTypeMfdDataImages) && key.HasValidLanguage() {
		holder := bitmaps.mfdArt[key.Language.ToIndex()].Get(res.ResourceID(key.Type))
		insertIndex := key.Index
		available := holder.BlockCount()

		if insertIndex >= available {
			insertIndex = available
		}
		writer := bytes.NewBuffer(nil)
		image.Write(writer, bmp, image.CompressedBitmap, true, 0)
		holder.SetBlockData(insertIndex, writer.Bytes())
		resultKey = model.MakeLocalizedResourceKey(key.Type, key.Language, insertIndex)
	} else {
		err = fmt.Errorf("Unsupported resource key %v", key)
	}

	return
}
