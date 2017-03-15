package model

import (
	"fmt"

	"github.com/inkyblackness/shocked-model"
)

// TextureAdapter is the entry point for game textures.
type TextureAdapter struct {
	context projectContext
	store   model.DataStore

	worldTextures               map[model.TextureSize]*Bitmaps
	worldTextureRequestsPending map[int]int
}

func newTextureAdapter(context projectContext, store model.DataStore) *TextureAdapter {
	adapter := &TextureAdapter{
		context: context,
		store:   store,

		worldTextures:               make(map[model.TextureSize]*Bitmaps),
		worldTextureRequestsPending: make(map[int]int)}

	for _, size := range model.TextureSizes() {
		adapter.worldTextures[size] = newBitmaps()
	}

	return adapter
}

func (adapter *TextureAdapter) clear() {
	for _, bitmaps := range adapter.worldTextures {
		bitmaps.clear()
	}
}

// WorldTextureCount returns the number of available textures.
func (adapter *TextureAdapter) WorldTextureCount() int {
	return 273
}

// RequestWorldTextureBitmaps will load the bitmap data for given world texture.
func (adapter *TextureAdapter) RequestWorldTextureBitmaps(key int) {
	if adapter.worldTextureRequestsPending[key] == 0 {
		for _, size := range model.TextureSizes() {
			adapter.requestWorldTextureBitmapInSize(key, size)
		}
	}
}

func (adapter *TextureAdapter) requestWorldTextureBitmapInSize(key int, size model.TextureSize) {
	adapter.worldTextureRequestsPending[key]++
	adapter.store.TextureBitmap(adapter.context.ActiveProjectID(), key, string(size),
		func(bmp *model.RawBitmap) {
			adapter.worldTextureRequestsPending[key]--
			adapter.worldTextures[size].setRawBitmap(key, bmp)
		},
		func() {
			adapter.worldTextureRequestsPending[key]--
			adapter.context.simpleStoreFailure(fmt.Sprintf("WorldTexture[%v][%v]", size, key))()
		})
}

// WorldTextures returns the container of bitmaps for given size.
func (adapter *TextureAdapter) WorldTextures(size model.TextureSize) *Bitmaps {
	return adapter.worldTextures[size]
}
