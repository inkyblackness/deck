package levelobj

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/data/interpreters"
)

// ForRealWorld returns an interpreter instance that handles the level class
// data of the specified object - in real world.
func ForRealWorld(objID res.ObjectID, data []byte) *interpreters.Instance {
	return realWorldEntries.specialize(int(objID.Class)).specialize(int(objID.Subclass)).specialize(int(objID.Type)).instance(data)
}
