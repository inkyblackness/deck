package levelobj

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/data/interpreters"
)

// ForCyberspace returns an interpreter instance that handles the level class
// data of the specified object - in cyberspace.
func ForCyberspace(objID res.ObjectID, data []byte) *interpreters.Instance {
	return cyberspaceEntries.specialize(int(objID.Class)).specialize(int(objID.Subclass)).specialize(int(objID.Type)).instance(data)
}
