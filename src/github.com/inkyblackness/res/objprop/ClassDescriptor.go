package objprop

// ClassDescriptor describes a single object class.
type ClassDescriptor struct {
	// GenericDataLength specifies the length of one generic type entry.
	GenericDataLength uint32
	// Subclasses contains descriptions of the subclasses of this class.
	// The index into the array is the subclass ID.
	Subclasses []SubclassDescriptor
}

// TotalDataLength returns the total length the class requires
func (desc ClassDescriptor) TotalDataLength() uint32 {
	total := uint32(0)

	for _, subclass := range desc.Subclasses {
		total += (desc.GenericDataLength * subclass.TypeCount) + subclass.TotalDataLength()
	}

	return total
}

// TotalTypeCount returns the total number of types in this class
func (desc ClassDescriptor) TotalTypeCount() uint32 {
	total := uint32(0)

	for _, subclass := range desc.Subclasses {
		total += subclass.TypeCount
	}

	return total
}
