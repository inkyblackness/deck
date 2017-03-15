package controls

import (
	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/ui"
)

// BitmapTexturizer creates a bitmap texture from bitmap information.
type BitmapTexturizer func(*graphics.Bitmap) *graphics.BitmapTexture

// Label is a control for displaying a text within an area.
type Label struct {
	area *ui.Area

	textPainter     graphics.TextPainter
	texturizer      BitmapTexturizer
	textureRenderer graphics.TextureRenderer

	scale             float32
	horizontalAligner Aligner
	verticalAligner   Aligner

	bitmap  graphics.TextBitmap
	texture *graphics.BitmapTexture
}

// Dispose releases all resources and removes the area from the tree.
func (label *Label) Dispose() {
	label.area.Remove()
	if label.texture != nil {
		label.texture.Dispose()
		label.texture = nil
	}
}

// SetText updates the current label text.
func (label *Label) SetText(text string) {
	if label.texture != nil {
		label.texture.Dispose()
		label.texture = nil
	}

	label.bitmap = label.textPainter.Paint(text)
	label.texture = label.texturizer(&label.bitmap.Bitmap)
}

func (label *Label) onRender(area *ui.Area) {
	u, v := label.texture.UV()
	fromLeft := float32(0.0)
	fromTop := float32(0.0)
	fromRight := u
	fromBottom := v
	textWidth, textHeight := label.texture.Size()
	scaledWidth, scaledHeight := textWidth*label.scale, textHeight*label.scale
	areaLeft := area.Left().Value()
	areaRight := area.Right().Value()
	areaWidth := areaRight - areaLeft
	areaTop := area.Top().Value()
	areaBottom := area.Bottom().Value()
	areaHeight := areaBottom - areaTop

	toLeft := areaLeft + label.horizontalAligner(areaWidth, scaledWidth)
	toTop := areaTop + label.verticalAligner(areaHeight, scaledHeight)
	toRight := toLeft + scaledWidth
	toBottom := toTop + scaledHeight

	if toLeft < areaLeft {
		fromLeft += (u / textWidth) * (areaLeft - toLeft) / label.scale
		toLeft = areaLeft
	}
	if toRight > areaRight {
		fromRight -= (u / textWidth) * (toRight - areaRight) / label.scale
		toRight = areaRight
	}
	if toTop < areaTop {
		fromTop += (v / textHeight) * (areaTop - toTop) / label.scale
		toTop = areaTop
	}
	if toBottom > areaBottom {
		fromBottom -= (v / textHeight) * (toBottom - areaBottom) / label.scale
		toBottom = areaBottom
	}

	modelMatrix := mgl.Ident4().Mul4(mgl.Translate3D(toLeft, toTop, 0.0)).Mul4(mgl.Scale3D(toRight-toLeft, toBottom-toTop, 1.0))

	label.textureRenderer.Render(&modelMatrix, label.texture, graphics.RectByCoord(fromLeft, fromTop, fromRight, fromBottom))
}
