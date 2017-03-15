package graphics

// TextureKey is a reference identifier for textures.
type TextureKey int

// ToInt returns the integer representation of the key.
func (key TextureKey) ToInt() int {
	return int(key)
}

// TextureKeyFromInt wraps an integer as TextureKey.
func TextureKeyFromInt(value int) TextureKey {
	return TextureKey(value)
}

// TextureQuery is a callback function to request the data for a specific texture.
type TextureQuery func(id TextureKey)

// BufferedTextureStore keeps textures in a buffer.
type BufferedTextureStore struct {
	query    TextureQuery
	textures map[TextureKey]*BitmapTexture
}

// NewBufferedTextureStore returns a new instance of a store.
func NewBufferedTextureStore(query TextureQuery) *BufferedTextureStore {
	return &BufferedTextureStore{
		query:    query,
		textures: make(map[TextureKey]*BitmapTexture)}
}

// Reset clears the store. It disposes any registered texture.
func (store *BufferedTextureStore) Reset() {
	oldTextures := store.textures

	store.textures = make(map[TextureKey]*BitmapTexture)
	for _, texture := range oldTextures {
		if texture != nil {
			texture.Dispose()
		}
	}
}

// Texture returns the texture associated with the given ID. May be nil if
// not yet known/available.
func (store *BufferedTextureStore) Texture(id TextureKey) *BitmapTexture {
	texture, existing := store.textures[id]

	if !existing {
		store.textures[id] = nil
		store.query(id)
	}

	return texture
}

// SetTexture registers a (new) texture under given ID. It disposes any
// previously registered texture.
func (store *BufferedTextureStore) SetTexture(id TextureKey, texture *BitmapTexture) {
	oldTexture := store.textures[id]

	store.textures[id] = texture
	if oldTexture != nil {
		oldTexture.Dispose()
	}
}
