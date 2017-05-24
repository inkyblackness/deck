package gameobj

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/data/interpreters"
)

// SpecificProperties returns an interpreter specific for the given object class and subclass.
func SpecificProperties(objID res.ObjectID, data []byte) *interpreters.Instance {
	desc := specificDescriptions[objID]
	if desc == nil {
		desc = specificDescriptions[res.MakeObjectID(objID.Class, objID.Subclass, anyType)]
	}
	if desc == nil {
		desc = interpreters.New()
	}

	return desc.For(data)
}
