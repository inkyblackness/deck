package objprop

import (
	"github.com/inkyblackness/res"
)

// ObjectIDToIndex returns the linear index for a given object ID based on a set of class descriptors.
func ObjectIDToIndex(desc []ClassDescriptor, id res.ObjectID) int {
	index := -1
	intClass := int(id.Class)

	if intClass < len(desc) {
		temp := uint32(0)

		for i := 0; i < intClass; i++ {
			temp += desc[i].TotalTypeCount()
		}

		classDesc := desc[intClass]
		intSubclass := int(id.Subclass)
		if intSubclass < len(classDesc.Subclasses) {
			for i := 0; i < intSubclass; i++ {
				temp += classDesc.Subclasses[i].TypeCount
			}

			subclassDesc := classDesc.Subclasses[intSubclass]
			intType := uint32(id.Type)
			if intType < subclassDesc.TypeCount {
				index = int(temp + intType)
			}
		}
	}

	return index
}
