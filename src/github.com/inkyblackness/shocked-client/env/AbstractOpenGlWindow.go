package env

// AbstractOpenGlWindow implements the common, basic functionality of OpenGlWindow
type AbstractOpenGlWindow struct {
	CallRender            RenderCallback
	CallOnMouseMove       MouseMoveCallback
	CallOnMouseButtonUp   MouseButtonCallback
	CallOnMouseButtonDown MouseButtonCallback
	CallOnMouseScroll     MouseScrollCallback
}

// InitAbstractOpenGlWindow returns an initialized instance.
func InitAbstractOpenGlWindow() AbstractOpenGlWindow {
	return AbstractOpenGlWindow{
		CallRender:            func() {},
		CallOnMouseMove:       func(float32, float32) {},
		CallOnMouseButtonUp:   func(uint32, uint32) {},
		CallOnMouseButtonDown: func(uint32, uint32) {}}
}

// OnRender implements the OpenGlWindow interface.
func (window *AbstractOpenGlWindow) OnRender(callback RenderCallback) {
	window.CallRender = callback
}

// OnMouseMove implements the OpenGlWindow interface.
func (window *AbstractOpenGlWindow) OnMouseMove(callback MouseMoveCallback) {
	window.CallOnMouseMove = callback
}

// OnMouseButtonDown implements the OpenGlWindow interface.
func (window *AbstractOpenGlWindow) OnMouseButtonDown(callback MouseButtonCallback) {
	window.CallOnMouseButtonDown = callback
}

// OnMouseButtonUp implements the OpenGlWindow interface.
func (window *AbstractOpenGlWindow) OnMouseButtonUp(callback MouseButtonCallback) {
	window.CallOnMouseButtonUp = callback
}

// OnMouseScroll implements the OpenGlWindow interface.
func (window *AbstractOpenGlWindow) OnMouseScroll(callback MouseScrollCallback) {
	window.CallOnMouseScroll = callback
}
