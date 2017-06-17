package model

// ResourceType is an enumeration of resource clusters.
type ResourceType uint16

const (
	// ResourceTypeMfdDataImages refers to the bitmaps used in the MFD data displays, such as logs.
	ResourceTypeMfdDataImages = ResourceType(0x0028)
)
