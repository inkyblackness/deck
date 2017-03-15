package model

import (
	"fmt"
)

// ObjectID is the reference of a specific game object.
type ObjectID int

// MakeObjectID returns a combined object identifier.
func MakeObjectID(class, subclass, objType int) ObjectID {
	return ObjectID((class << 16) | (subclass << 8) | (objType))
}

// ObjectIDFromInt returns an object identifier wrapping the provided integer.
func ObjectIDFromInt(value int) ObjectID {
	return ObjectID(value)
}

// ToInt returns a single integer representation of the ID.
func (id ObjectID) ToInt() int {
	return int(id)
}

// Class returns the class value.
func (id ObjectID) Class() int {
	return (id.ToInt() >> 16) & 0xFF
}

// Subclass returns the subclass value.
func (id ObjectID) Subclass() int {
	return (id.ToInt() >> 8) & 0xFF
}

// Type returns the type value.
func (id ObjectID) Type() int {
	return (id.ToInt() >> 0) & 0xFF
}

// String implements the Stringer interface.
func (id ObjectID) String() string {
	return fmt.Sprintf("%2d/%d/%2d", id.Class(), id.Subclass(), id.Type())
}
