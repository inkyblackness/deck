package controls

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/ui"
)

// ImageProvider returns the image to render in the display.
type ImageProvider func() *graphics.BitmapTexture

// ImageDisplay is a graphical selection tool for textures.
type ImageDisplay struct {
	area *ui.Area

	textureRenderer graphics.TextureRenderer
	provider        ImageProvider
}

// Dispose releases all current resources.
func (control *ImageDisplay) Dispose() {
	control.area.Remove()
}

func (control *ImageDisplay) onRender(area *ui.Area) {
	image := control.provider()

	if image != nil {
		areaTop := area.Top().Value()
		areaBottom := area.Bottom().Value()
		areaLeft := area.Left().Value()
		areaRight := area.Right().Value()
		areaHeight := areaBottom - areaTop
		areaWidth := areaRight - areaLeft
		fromLeft := float32(0.0)
		fromTop := float32(0.0)
		fromRight, fromBottom := image.UV()
		imageWidth, imageHeight := image.Size()
		widthFitting := areaWidth / imageWidth
		heightFitting := areaHeight / imageHeight
		scale := float32(math.Min(float64(widthFitting), float64(heightFitting)))
		toLeft := areaLeft + (areaWidth-(imageWidth*scale))/2.0
		toTop := areaTop + (areaHeight-(imageHeight*scale))/2.0

		modelMatrix := mgl.Ident4().Mul4(mgl.Translate3D(toLeft, toTop, 0)).Mul4(mgl.Scale3D(imageWidth*scale, imageHeight*scale, 1.0))
		control.textureRenderer.Render(&modelMatrix, image, graphics.RectByCoord(fromLeft, fromTop, fromRight, fromBottom))
	}
}
