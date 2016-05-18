package textprop

// Provider wraps the Provide method.
type Provider interface {
	// EntryCount returns the amount of entries available
	EntryCount() uint32
	// Provide returns the entry data for the requested id.
	Provide(id uint32) []byte
}
