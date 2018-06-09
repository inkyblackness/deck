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

	gameTextures *observable
}

func newTextureAdapter(context projectContext, store model.DataStore) *TextureAdapter {
	adapter := &TextureAdapter{
		context: context,
		store:   store,

		worldTextures:               make(map[model.TextureSize]*Bitmaps),
		worldTextureRequestsPending: make(map[int]int),

		gameTextures: newObservable()}

	for _, size := range model.TextureSizes() {
		adapter.worldTextures[size] = newBitmaps()
	}
	adapter.gameTextures.set(&[]*GameTexture{})

	return adapter
}

func (adapter *TextureAdapter) clear() {
	for _, bitmaps := range adapter.worldTextures {
		bitmaps.clear()
	}
	adapter.gameTextures.set(&[]*GameTexture{})
}

func (adapter *TextureAdapter) refresh() {
	adapter.store.Textures(adapter.context.ActiveProjectID(), adapter.onNewGameTextures, adapter.context.simpleStoreFailure("Textures"))
}

func (adapter *TextureAdapter) onNewGameTextures(properties []model.TextureProperties) {
	count := len(properties)
	textures := make([]*GameTexture, count)

	for id := 0; id < count; id++ {
		texture := newGameTexture(id)
		texture.properties = properties[id]
		textures[id] = texture
	}

	adapter.gameTextures.set(&textures)
}

// OnGameTexturesChanged registers a callback for updates.
func (adapter *TextureAdapter) OnGameTexturesChanged(callback func()) {
	adapter.gameTextures.addObserver(callback)
}

func (adapter *TextureAdapter) gameTextureList() []*GameTexture {
	return *adapter.gameTextures.get().(*[]*GameTexture)
}

// WorldTextureCount returns the number of available textures.
func (adapter *TextureAdapter) WorldTextureCount() int {
	return len(adapter.gameTextureList())
}

// GameTexture returns texture information for identified texture.
func (adapter *TextureAdapter) GameTexture(id int) (texture *GameTexture) {
	list := adapter.gameTextureList()
	if (id >= 0) && (id < len(list)) {
		texture = list[id]
	} else {
		texture = nullGameTexture(id)
	}
	return
}

// RequestTexturePropertiesChange requests to change properties of a single texture.
func (adapter *TextureAdapter) RequestTexturePropertiesChange(id int, properties *model.TextureProperties) {
	textures := adapter.gameTextureList()

	adapter.store.SetTextureProperties(adapter.context.ActiveProjectID(), id, properties,
		func(updatedProperties *model.TextureProperties) {
			textures[id].properties = *updatedProperties
			adapter.gameTextures.notifyObservers()
		}, adapter.context.simpleStoreFailure("SetTextureProperties"))
}

// TextureBitmap returns the raw bitmap for given key - if available.
func (adapter *TextureAdapter) TextureBitmap(id int, size model.TextureSize) *model.RawBitmap {
	return adapter.worldTextures[size].RawBitmap(id)
}

// RequestTextureBitmapChange requests to change the bitmap of a single texture.
func (adapter *TextureAdapter) RequestTextureBitmapChange(id int, size model.TextureSize, rawBitmap *model.RawBitmap) {
	adapter.store.SetTextureBitmap(adapter.context.ActiveProjectID(), id, string(size), rawBitmap,
		func(rawResult *model.RawBitmap) {
			adapter.worldTextures[size].setRawBitmap(id, rawResult)
		}, adapter.context.simpleStoreFailure("SetTextureBitmap"))
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
