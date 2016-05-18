package textprop

// Store represents a dynamically accessible container of texture properties.
type Store interface {
	// Get returns the data for the requested ID.
	Get(id uint32) []byte

	// Put takes the provided data and associates it with the given ID.
	Put(id uint32, data []byte)
}
