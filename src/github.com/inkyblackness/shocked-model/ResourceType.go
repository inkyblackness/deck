package model

// ResourceType is an enumeration of resource clusters.
type ResourceType uint16

const (
	// ResourceTypeMfdDataImages refers to the bitmaps used in the MFD data displays, such as logs.
	ResourceTypeMfdDataImages = ResourceType(0x0028)

	// ResourceTypeTrapMessages refers to the texts shown in the MFD based on actions.
	ResourceTypeTrapMessages = ResourceType(0x0867)
	// ResourceTypeWords refers to the texts of WORDS level objects.
	ResourceTypeWords = ResourceType(0x0868)
	// ResourceTypeLogCategories contains the category names of logs.
	ResourceTypeLogCategories = ResourceType(0x0870)
	// ResourceTypeVariousMessages contains all sorts of messages, including door lock messages.
	ResourceTypeVariousMessages = ResourceType(0x0871)
	// ResourceTypeScreenMessages contains the messages shown on screens.
	ResourceTypeScreenMessages = ResourceType(0x0877)
	// ResourceTypeInfoNodeMessages contains the short text fragments found in cyberspace of 8/5/6 nodes.
	ResourceTypeInfoNodeMessages = ResourceType(0x0878)
	// ResourceTypeAccessCardNames contains the names of the access cards.
	ResourceTypeAccessCardNames = ResourceType(0x0879)
	// ResourceTypeDataletMessages contains the short text fragments found in cyberspace of 8/5/8 nodes.
	ResourceTypeDataletMessages = ResourceType(0x087A)

	// ResourceTypePaperTexts refers to the texts written on lose papers.
	ResourceTypePaperTexts = ResourceType(0x003C)

	// ResourceTypeTrapAudio refers to the audio played along trap messages.
	ResourceTypeTrapAudio = ResourceType(0x0C1C)
)
