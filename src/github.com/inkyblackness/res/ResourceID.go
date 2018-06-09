package res

// ResourceID uniquely identifies a resource
type ResourceID uint16

// Value returns the identifiers raw integer value.
func (id ResourceID) Value() uint16 {
	return uint16(id)
}
