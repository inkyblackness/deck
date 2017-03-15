package display

import (
	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-client/editor/model"
	"github.com/inkyblackness/shocked-client/graphics"
)

// Context provides some global resources.
type Context interface {
	ModelAdapter() *model.Adapter
	NewRenderContext(viewMatrix *mgl.Mat4) *graphics.RenderContext
	ForGraphics() graphics.Context
}
