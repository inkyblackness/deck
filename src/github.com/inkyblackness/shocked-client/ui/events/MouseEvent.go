package events

// MouseEvent is the base structure for all mouse related events.
type MouseEvent struct {
	eventType EventType

	x, y     float32
	modifier uint32
	buttons  uint32
}

// InitMouseEvent initializes a basic mouse event structure.
func InitMouseEvent(eventType EventType, x, y float32, modifier uint32, buttons uint32) MouseEvent {
	event := MouseEvent{
		eventType: eventType,
		x:         x,
		y:         y,
		modifier:  modifier,
		buttons:   buttons}

	return event
}

// EventType implements the Event interface.
func (event *MouseEvent) EventType() EventType {
	return event.eventType
}

// Position returns the coordinate of the event.
func (event *MouseEvent) Position() (x, y float32) {
	return event.x, event.y
}

// Modifier returns a bitmask of currently active keyboard modifier.
func (event *MouseEvent) Modifier() uint32 {
	return event.modifier
}

// Buttons returns a bitmask of currently pressed mouse buttons.
func (event *MouseEvent) Buttons() uint32 {
	return event.buttons
}
