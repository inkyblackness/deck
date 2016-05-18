package graphics

// Texture describes a texture in graphics memory.
type Texture interface {
	// Dispose releases any internal resources.
	Dispose()
	// Handle returns the texture handle.
	Handle() uint32
}
