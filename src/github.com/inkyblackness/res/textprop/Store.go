package textprop

// Store represents a dynamically accessible container of texture properties.
type Store interface {
	// EntryCount returns the number of available textures.
	EntryCount() uint32

	// Get returns the data for the requested ID.
	Get(id uint32) []byte

	// Put takes the provided data and associates it with the given ID.
	Put(id uint32, data []byte)
}
