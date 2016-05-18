package objprop

const (
	// CommonPropertiesLength specifies the length a common properties structure has.
	CommonPropertiesLength = uint32(27)
)

// StandardProperties returns an array of class descriptors that represent the standard
// configuration of the existing file
func StandardProperties() []ClassDescriptor {
	result := []ClassDescriptor{}

	{ // Weapons
		subclasses := []SubclassDescriptor{
			SubclassDescriptor{TypeCount: 5, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 2, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 2, SpecificDataLength: 16},
			SubclassDescriptor{TypeCount: 2, SpecificDataLength: 13},
			SubclassDescriptor{TypeCount: 3, SpecificDataLength: 13},
			SubclassDescriptor{TypeCount: 2, SpecificDataLength: 18}}

		result = append(result, ClassDescriptor{GenericDataLength: 2, Subclasses: subclasses})
	}
	{ // Clips
		subclasses := []SubclassDescriptor{
			SubclassDescriptor{TypeCount: 2, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 2, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 3, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 2, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 2, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 2, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 2, SpecificDataLength: 1}}

		result = append(result, ClassDescriptor{GenericDataLength: 14, Subclasses: subclasses})
	}
	{ // Projectiles
		subclasses := []SubclassDescriptor{
			SubclassDescriptor{TypeCount: 6, SpecificDataLength: 20},
			SubclassDescriptor{TypeCount: 16, SpecificDataLength: 6},
			SubclassDescriptor{TypeCount: 2, SpecificDataLength: 1}}

		result = append(result, ClassDescriptor{GenericDataLength: 1, Subclasses: subclasses})
	}
	{ // Explosives
		subclasses := []SubclassDescriptor{
			SubclassDescriptor{TypeCount: 5, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 3, SpecificDataLength: 3}}

		result = append(result, ClassDescriptor{GenericDataLength: 15, Subclasses: subclasses})
	}
	{ // Patches
		subclasses := []SubclassDescriptor{
			SubclassDescriptor{TypeCount: 7, SpecificDataLength: 1}}

		result = append(result, ClassDescriptor{GenericDataLength: 22, Subclasses: subclasses})
	}
	{ // Hardware
		subclasses := []SubclassDescriptor{
			SubclassDescriptor{TypeCount: 5, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 10, SpecificDataLength: 1}}

		result = append(result, ClassDescriptor{GenericDataLength: 9, Subclasses: subclasses})
	}
	{ // Softs
		subclasses := []SubclassDescriptor{
			SubclassDescriptor{TypeCount: 7, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 3, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 4, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 5, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 3, SpecificDataLength: 1}}

		result = append(result, ClassDescriptor{GenericDataLength: 5, Subclasses: subclasses})
	}
	{ // Fixtures
		subclasses := []SubclassDescriptor{
			SubclassDescriptor{TypeCount: 9, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 10, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 11, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 4, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 9, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 8, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 16, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 10, SpecificDataLength: 1}}

		result = append(result, ClassDescriptor{GenericDataLength: 2, Subclasses: subclasses})
	}
	{ // Items
		subclasses := []SubclassDescriptor{
			SubclassDescriptor{TypeCount: 8, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 10, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 15, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 6, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 12, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 12, SpecificDataLength: 6},
			SubclassDescriptor{TypeCount: 9, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 8, SpecificDataLength: 2}}

		result = append(result, ClassDescriptor{GenericDataLength: 2, Subclasses: subclasses})
	}
	{ // Panels
		subclasses := []SubclassDescriptor{
			SubclassDescriptor{TypeCount: 9, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 7, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 3, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 11, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 2, SpecificDataLength: 0}, // skipped vending machines
			SubclassDescriptor{TypeCount: 3, SpecificDataLength: 1}}

		result = append(result, ClassDescriptor{GenericDataLength: 1, Subclasses: subclasses})
	}
	{ // Barriers
		subclasses := []SubclassDescriptor{
			SubclassDescriptor{TypeCount: 10, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 9, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 7, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 5, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 10, SpecificDataLength: 1}}

		result = append(result, ClassDescriptor{GenericDataLength: 1, Subclasses: subclasses})
	}
	{ // Animated
		subclasses := []SubclassDescriptor{
			SubclassDescriptor{TypeCount: 9, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 11, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 14, SpecificDataLength: 1}}

		result = append(result, ClassDescriptor{GenericDataLength: 2, Subclasses: subclasses})
	}
	{ // Marker
		subclasses := []SubclassDescriptor{
			SubclassDescriptor{TypeCount: 13, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 1, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 5, SpecificDataLength: 1}}

		result = append(result, ClassDescriptor{GenericDataLength: 1, Subclasses: subclasses})
	}
	{ // Container
		subclasses := []SubclassDescriptor{
			SubclassDescriptor{TypeCount: 3, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 3, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 4, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 8, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 13, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 7, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 8, SpecificDataLength: 1}}

		result = append(result, ClassDescriptor{GenericDataLength: 3, Subclasses: subclasses})
	}
	{ // Critter
		subclasses := []SubclassDescriptor{
			SubclassDescriptor{TypeCount: 9, SpecificDataLength: 1},
			SubclassDescriptor{TypeCount: 12, SpecificDataLength: 2},
			SubclassDescriptor{TypeCount: 7, SpecificDataLength: 2},
			SubclassDescriptor{TypeCount: 7, SpecificDataLength: 6},
			SubclassDescriptor{TypeCount: 2, SpecificDataLength: 1}}

		result = append(result, ClassDescriptor{GenericDataLength: 75, Subclasses: subclasses})
	}

	return result
}
