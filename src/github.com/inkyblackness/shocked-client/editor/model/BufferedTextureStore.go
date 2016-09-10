package model

import (
	"fmt"

	"github.com/inkyblackness/shocked-client/graphics"
)

// TextureKey is a reference identifier for textures.
// It has to implement fmt.Stringer to uniquely represent itself.
type TextureKey interface {
	fmt.Stringer
}

// TextureQuery is a callback function to request the data for a specific texture.
type TextureQuery func(id TextureKey)

// BufferedTextureStore keeps textures in a buffer.
type BufferedTextureStore struct {
	query    TextureQuery
	textures map[string]graphics.Texture
}

// NewBufferedTextureStore returns a new instance of a store.
func NewBufferedTextureStore(query TextureQuery) *BufferedTextureStore {
	return &BufferedTextureStore{
		query:    query,
		textures: make(map[string]graphics.Texture)}
}

// Reset clears the store. It disposes any registered texture.
func (store *BufferedTextureStore) Reset() {
	oldTextures := store.textures

	store.textures = make(map[string]graphics.Texture)
	for _, texture := range oldTextures {
		texture.Dispose()
	}
}

// Texture returns the texture associated with the given ID. May be null if
// not yet known/available.
func (store *BufferedTextureStore) Texture(id TextureKey) graphics.Texture {
	idString := id.String()
	texture, existing := store.textures[idString]

	if !existing {
		store.query(id)
		store.textures[idString] = nil
	}

	return texture
}

// SetTexture registers a (new) texture under given ID. It disposes any
// previously registered texture.
func (store *BufferedTextureStore) SetTexture(id TextureKey, texture graphics.Texture) {
	idString := id.String()
	oldTexture := store.textures[idString]

	store.textures[idString] = texture
	if oldTexture != nil {
		oldTexture.Dispose()
	}
}
