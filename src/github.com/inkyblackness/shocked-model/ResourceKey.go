package model

import (
	"fmt"
)

// ResourceKey is the reference of a specific game resource.
type ResourceKey struct {
	Type     ResourceType
	Language ResourceLanguage
	Index    uint16
}

// MakeResourceKey returns a combined resource identifier.
func MakeResourceKey(resourceType ResourceType, index uint16) ResourceKey {
	return ResourceKey{Type: resourceType, Language: ResourceLanguageUnspecific, Index: index}
}

// MakeLocalizedResourceKey returns a combined resource identifier with specified language.
func MakeLocalizedResourceKey(resourceType ResourceType, language ResourceLanguage, index uint16) ResourceKey {
	return ResourceKey{Type: resourceType, Language: language, Index: index}
}

// ResourceKeyFromInt returns a resource identifier wrapping the provided integer.
func ResourceKeyFromInt(value int) ResourceKey {
	return ResourceKey{
		Type:     ResourceType(uint16((value >> 16) & 0xFFFF)),
		Language: ResourceLanguage((value >> 13) & 0x3),
		Index:    uint16(value & 0x1FFF)}
}

// HasValidLanguage returns true if the Language field is within range [1..3].
func (id ResourceKey) HasValidLanguage() bool {
	return (int(id.Language) >= 1) && (int(id.Language) <= 3)
}

// ToInt returns a single integer representation of the ID.
func (id ResourceKey) ToInt() int {
	return (int(id.Type) << 16) | (int(id.Language) << 13) | int(id.Index)
}

// String implements the Stringer interface.
func (id ResourceKey) String() string {
	languages := []string{"*", "STD", "FRN", "GER"}

	return fmt.Sprintf("0x%04X:%03d[%v]", uint16(id.Type), id.Index, languages[id.Language])
}
