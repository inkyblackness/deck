package objprop

// ObjectData holds object-type specific data
type ObjectData struct {
	// Class-specific data
	Generic []byte
	// Type-specific data
	Specific []byte
	// Common properties (format shared between all objects)
	Common []byte
}
