package events

// ClipboardCopyEventType is the name for events where data shall be written to the clipboard.
const ClipboardCopyEventType = EventType("clipboard.copy")

// ClipboardPasteEventType is the name for events where data shall be read from the clipboard.
const ClipboardPasteEventType = EventType("clipboard.paste")

// ClipboardEvent is the base structure for all clipboard related events.
type ClipboardEvent struct {
	eventType EventType

	x, y      float32
	clipboard Clipboard
}

// NewClipboardEvent initializes a basic clipboard event structure.
func NewClipboardEvent(eventType EventType, x, y float32, clipboard Clipboard) *ClipboardEvent {
	event := &ClipboardEvent{
		eventType: eventType,
		x:         x,
		y:         y,
		clipboard: clipboard}

	return event
}

// EventType implements the Event interface.
func (event *ClipboardEvent) EventType() EventType {
	return event.eventType
}

// Position returns the coordinate of the event.
func (event *ClipboardEvent) Position() (x, y float32) {
	return event.x, event.y
}

// Clipboard returns access to the clipboard.
func (event *ClipboardEvent) Clipboard() Clipboard {
	return event.clipboard
}
