package events

// MouseScrollEvent describes a scrolling operation.
type MouseScrollEvent struct {
	MouseEvent

	dx, dy float32
}

// MouseScrollEventType is the name for events where the mouse scroll function was used.
const MouseScrollEventType = EventType("mouse.scroll")

// NewMouseScrollEvent returns a new instance of a mouse event.
func NewMouseScrollEvent(x, y float32, modifier uint32, buttons uint32, dx, dy float32) *MouseScrollEvent {
	event := &MouseScrollEvent{
		MouseEvent: InitMouseEvent(MouseScrollEventType, x, y, modifier, buttons),
		dx:         dx,
		dy:         dy}

	return event
}

// Deltas returns the offsets of the scrolling operation.
func (event *MouseScrollEvent) Deltas() (dx, dy float32) {
	return event.dx, event.dy
}
