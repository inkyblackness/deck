package env

import (
	"github.com/inkyblackness/shocked-client/opengl"
)

// RenderCallback is the function to receive render events. When the callback
// returns, the window will swap the internal buffer.
type RenderCallback func()

// MouseMoveCallback is the function to receive the current mouse coordinate while moving.
// Movement is reported while the cursor is within the client area of the window, and
// beyond the window as long as at least one captured button is pressed.
// Reported values are with sub-pixel precision, if possible.
type MouseMoveCallback func(x float32, y float32)

// MouseButtonCallback is the function to receive button up/down events.
// An Up event is sent for every reported Down event, even if the mouse cursor is outside
// the client area.
type MouseButtonCallback func(buttonMask uint32, modifierMask uint32)

// MouseScrollCallback is the function to receive scroll events.
// Delta values are right-hand oriented: positive values go right/down/far.
type MouseScrollCallback func(dx float32, dy float32)

// OpenGlWindow represents an OpenGL render surface.
type OpenGlWindow interface {
	// OpenGl returns the OpenGL API wrapper for this window.
	OpenGl() opengl.OpenGl
	// OnRender registers a callback function which shall be called to update the scene.
	OnRender(callback RenderCallback)

	// Size returns the dimensions of the window display area in pixel.
	Size() (width int, height int)

	// OnMouseMove registers a callback function for mouse move events.
	OnMouseMove(callback MouseMoveCallback)
	// OnMouseButtonDown registers a callback function for mouse button down events.
	OnMouseButtonDown(callback MouseButtonCallback)
	// OnMouseButtonUp registers a callback function for mouse button up events.
	OnMouseButtonUp(callback MouseButtonCallback)
	// OnMouseScroll registers a callback function for mouse scroll events.
	OnMouseScroll(callback MouseScrollCallback)
}
