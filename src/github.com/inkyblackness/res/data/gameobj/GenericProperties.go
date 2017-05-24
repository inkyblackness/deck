package gameobj

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/data/interpreters"
)

// GenericProperties returns an interpreter specific for the given object class.
func GenericProperties(objClass res.ObjectClass, data []byte) *interpreters.Instance {
	desc := genericDescriptions[objClass]
	if desc == nil {
		desc = interpreters.New()
	}

	return desc.For(data)
}
