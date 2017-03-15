package events

// EventType is an identifier of events.
type EventType string

// Event is the root interface for all UI events. Its EventType property
// specifies what kind of event it is.
type Event interface {
	// EventType returns the identifier for the event.
	EventType() EventType
}
