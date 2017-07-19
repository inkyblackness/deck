package model

import (
	"fmt"
)

// ObjectBitmapID identifies a bitmap of an object
type ObjectBitmapID struct {
	ObjectID ObjectID
	Index    int
}

// ObjectBitmapIDFromInt returns an object bitmap identifier wrapping the provided integer.
func ObjectBitmapIDFromInt(value int) (id ObjectBitmapID) {
	id.ObjectID = ObjectID(value & 0x00FFFFFF)
	id.Index = value >> 24
	return
}

// ToInt returns a single integer representation of the ID.
func (id ObjectBitmapID) ToInt() int {
	return id.ObjectID.ToInt() + (id.Index << 24)
}

func (id ObjectBitmapID) String() string {
	return fmt.Sprintf("%v[%v]", id.ObjectID, id.Index)
}
