package model

import (
	"fmt"
)

// ObjectID is the reference of a specific game object.
type ObjectID struct {
	class    int
	subclass int
	objType  int
}

// MakeObjectID returns a combined object identifier.
func MakeObjectID(class, subclass, objType int) ObjectID {
	return ObjectID{class, subclass, objType}
}

// Class returns the class value.
func (id ObjectID) Class() int {
	return id.class
}

// Subclass returns the subclass value.
func (id ObjectID) Subclass() int {
	return id.subclass
}

// Type returns the type value.
func (id ObjectID) Type() int {
	return id.objType
}

// String implements the Stringer interface.
func (id ObjectID) String() string {
	return fmt.Sprintf("%2d/%d/%2d", id.class, id.subclass, id.objType)
}
