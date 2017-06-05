package conditions

import (
	"github.com/inkyblackness/res/data/interpreters"
)

var objectType = interpreters.New().
	With("ObjectType", 0, 3).As(interpreters.SpecialValue("ObjectType"))

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
