package events

// MouseMoveEvent describes an update of the position of the mouse.
type MouseMoveEvent struct {
	MouseEvent
}

// MouseMoveEventType is the type for mouse move events.
const MouseMoveEventType = EventType("mouse.move")

// NewMouseMoveEvent returns a new instance of a mouse event.
func NewMouseMoveEvent(x, y float32, modifier uint32, buttons uint32) *MouseMoveEvent {
	event := &MouseMoveEvent{
		MouseEvent: InitMouseEvent(MouseMoveEventType, x, y, modifier, buttons)}

	return event
}
