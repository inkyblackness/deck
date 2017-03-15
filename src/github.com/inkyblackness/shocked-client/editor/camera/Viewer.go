package camera

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

// Viewer represents the view of a camera.
type Viewer interface {
	// ViewMatrix returns the transformation matrix for the current view.
	ViewMatrix() *mgl.Mat4
}
