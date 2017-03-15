package events

// MouseButtonEvent describes an update of buttons of the mouse.
type MouseButtonEvent struct {
	MouseEvent

	affectedButtons uint32
}

// MouseButtonDownEventType is the name for events where buttons where pressed.
const MouseButtonDownEventType = EventType("mouse.button.down")

// MouseButtonUpEventType is the name for events where buttons where de-pressed.
const MouseButtonUpEventType = EventType("mouse.button.up")

// MouseButtonClickedEventType is the name for events where buttons where clicked without moving the cursor.
const MouseButtonClickedEventType = EventType("mouse.button.clicked")

// NewMouseButtonEvent returns a new instance of a mouse event.
func NewMouseButtonEvent(eventType EventType, x, y float32, modifier uint32, buttons uint32, affectedButtons uint32) *MouseButtonEvent {
	event := &MouseButtonEvent{
		MouseEvent:      InitMouseEvent(eventType, x, y, modifier, buttons),
		affectedButtons: affectedButtons}

	return event
}

// AffectedButtons returns the bitmask of the buttons that were affected by the event.
// The bitmask returned by Buttons() returns the state of the mouse buttons after this event.
func (event *MouseButtonEvent) AffectedButtons() uint32 {
	return event.affectedButtons
}
