package graphics

import (
	"github.com/inkyblackness/shocked-model"
)

// Context is a provider of graphic utilities.
type Context interface {
	RectangleRenderer() *RectangleRenderer
	TextPainter() TextPainter
	Texturize(bmp *Bitmap) *BitmapTexture
	UITextRenderer() *BitmapTextureRenderer

	NewPaletteTexture(colorProvider ColorProvider) *PaletteTexture
	BitmapsStore() *BufferedTextureStore
	WorldTextureStore(size model.TextureSize) *BufferedTextureStore
	GameObjectIconsStore() *BufferedTextureStore
}
