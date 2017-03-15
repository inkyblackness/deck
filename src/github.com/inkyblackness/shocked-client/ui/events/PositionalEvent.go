package events

// PositionalEvent is an event happening at a specific position.
type PositionalEvent interface {
	Event

	// Position returns the location of the event.
	Position() (float32, float32)
}
