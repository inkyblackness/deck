package graphics

import (
	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-client/opengl"
)

// RenderContext provides current render properties.
type RenderContext struct {
	gl opengl.OpenGl

	viewMatrix       *mgl.Mat4
	projectionMatrix *mgl.Mat4
}

// NewBasicRenderContext returns a render context for the provided parameters.
func NewBasicRenderContext(gl opengl.OpenGl, projectionMatrix *mgl.Mat4, viewMatrix *mgl.Mat4) *RenderContext {
	return &RenderContext{
		gl:               gl,
		viewMatrix:       viewMatrix,
		projectionMatrix: projectionMatrix}
}

// OpenGl returns the GL interface.
func (context *RenderContext) OpenGl() opengl.OpenGl {
	return context.gl
}

// ViewMatrix returns the current view matrix.
func (context *RenderContext) ViewMatrix() *mgl.Mat4 {
	return context.viewMatrix
}

// ProjectionMatrix returns the current projection matrix.
func (context *RenderContext) ProjectionMatrix() *mgl.Mat4 {
	return context.projectionMatrix
}
