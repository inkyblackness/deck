package objprop

// SubclassDescriptor describes one subclass.
type SubclassDescriptor struct {
	// TypeCount specifies how many types exist in this subclass.
	TypeCount uint32
	// SpecificDataLength specifies the length of one specific type entry.
	SpecificDataLength uint32
}

// TotalDataLength returns the total length the subclass requires in the properties file
func (desc SubclassDescriptor) TotalDataLength() uint32 {
	return (desc.SpecificDataLength + CommonPropertiesLength) * desc.TypeCount
}
