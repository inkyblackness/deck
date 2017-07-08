package model

// Maximum values for text resource types
const (
	MaxTrapMessages     = 512
	MaxWords            = 512
	MaxLogCategories    = 16
	MaxScreenMessages   = 120
	MaxInfoNodeMessages = 256
	MaxAccessCardNames  = 32 * 2
	MaxDataletMessages  = 256
	MaxPaperTexts       = 16
)

var maxEntriesByType = map[ResourceType]int{
	ResourceTypeTrapMessages:     MaxTrapMessages,
	ResourceTypeWords:            MaxWords,
	ResourceTypeLogCategories:    MaxLogCategories,
	ResourceTypeScreenMessages:   MaxScreenMessages,
	ResourceTypeInfoNodeMessages: MaxInfoNodeMessages,
	ResourceTypeAccessCardNames:  MaxAccessCardNames,
	ResourceTypeDataletMessages:  MaxDataletMessages,
	ResourceTypePaperTexts:       MaxPaperTexts}

// MaxEntriesFor returns the maximum count of resources of a given type.
func MaxEntriesFor(resourceType ResourceType) int {
	return maxEntriesByType[resourceType]
}
