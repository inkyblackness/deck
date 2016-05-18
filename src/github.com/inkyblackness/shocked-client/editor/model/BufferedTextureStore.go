package model

import (
	"github.com/inkyblackness/shocked-client/graphics"
)

// TextureQuery is a callback function to request the data for a specific texture.
type TextureQuery func(id int)

// BufferedTextureStore keeps textures in a buffer.
type BufferedTextureStore struct {
	query    TextureQuery
	textures map[int]graphics.Texture
}

// NewBufferedTextureStore returns a new instance of a store.
func NewBufferedTextureStore(query TextureQuery) *BufferedTextureStore {
	return &BufferedTextureStore{
		query:    query,
		textures: make(map[int]graphics.Texture)}
}

// Reset clears the store. It disposes any registered texture.
func (store *BufferedTextureStore) Reset() {
	oldTextures := store.textures

	store.textures = make(map[int]graphics.Texture)
	for _, texture := range oldTextures {
		texture.Dispose()
	}
}

// Texture returns the texture associated with the given ID. May be null if
// not yet known/available.
func (store *BufferedTextureStore) Texture(id int) graphics.Texture {
	texture, existing := store.textures[id]

	if !existing {
		store.query(id)
		store.textures[id] = nil
	}

	return texture
}

// SetTexture registers a (new) texture under given ID. It disposes any
// previously registered texture.
func (store *BufferedTextureStore) SetTexture(id int, texture graphics.Texture) {
	oldTexture := store.textures[id]

	store.textures[id] = texture
	if oldTexture != nil {
		oldTexture.Dispose()
	}
}
