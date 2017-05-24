package gameobj

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/data/interpreters"
)

var projectileGenerics = interpreters.New().
	With("Flags", 0, 1).As(interpreters.Bitfield(map[uint32]string{
	0x01: "EmitLight",
	0x02: "BounceOffWalls",
	0x04: "BouncePassObjects",
	0x08: "Unknown08"}))

var cyberProjectiles = interpreters.New().
	Refining("ColorScheme", 0, 6, cyberColorScheme, interpreters.Always)

func initProjectiles() {
	objClass := res.ObjectClass(2)

	genericDescriptions[objClass] = projectileGenerics

	setSpecificByType(objClass, 1, 9, cyberProjectiles)
	setSpecificByType(objClass, 1, 10, cyberProjectiles)
	setSpecificByType(objClass, 1, 11, cyberProjectiles)
	setSpecificByType(objClass, 1, 12, cyberProjectiles)
	setSpecificByType(objClass, 1, 13, cyberProjectiles)
}
