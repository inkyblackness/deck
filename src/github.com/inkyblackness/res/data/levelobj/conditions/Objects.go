package conditions

import (
	"github.com/inkyblackness/res/data/interpreters"
)

var objectType = interpreters.New().
	With("Type", 0, 1).
	With("Subclass", 1, 1).
	With("Class", 2, 1).As(interpreters.RangedValue(0, 14))

var objectIndex = interpreters.New().
	With("ObjectIndex", 0, 2).As(interpreters.ObjectIndex())

// ObjectType returns a condition description for object types.
func ObjectType() *interpreters.Description {
	return objectType
}

// ObjectIndex returns a condition description for object indices.
func ObjectIndex() *interpreters.Description {
	return objectIndex
}
