package display

import (
	"github.com/go-gl/mathgl/mgl32"
)

// RenderContext describes the current render properties
type RenderContext struct {
	viewportWidth  int
	viewportHeight int

	viewMatrix       mgl32.Mat4
	projectionMatrix mgl32.Mat4
}

// NewBasicRenderContext returns a render context for the provided parameters.
func NewBasicRenderContext(width, height int, viewMatrix mgl32.Mat4) *RenderContext {
	return &RenderContext{
		viewportWidth:    width,
		viewportHeight:   height,
		viewMatrix:       viewMatrix,
		projectionMatrix: mgl32.Ortho2D(0, float32(width), float32(height), 0)}
}

// ViewportSize returns the size of the current viewport, in pixel.
func (context *RenderContext) ViewportSize() (width int, height int) {
	return context.viewportWidth, context.viewportHeight
}

// ViewMatrix returns the current view matrix.
func (context *RenderContext) ViewMatrix() *mgl32.Mat4 {
	return &context.viewMatrix
}

// ProjectionMatrix returns the current projection matrix.
func (context *RenderContext) ProjectionMatrix() *mgl32.Mat4 {
	return &context.projectionMatrix
}
