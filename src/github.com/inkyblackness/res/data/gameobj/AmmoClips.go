package gameobj

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/data/interpreters"
)

var ammoClipGenerics = interpreters.New().
	Refining("BasicWeapon", 0, 8, basicWeapon, interpreters.Always).
	With("RoundsPerClip", 8, 1).
	With("ImpactForce", 9, 1).
	With("Kickback", 10, 2).As(interpreters.RangedValue(-10000, +10000)).
	With("Range", 12, 1).
	With("AimSkew", 13, 1)

func initAmmoClips() {
	objClass := res.ObjectClass(1)

	genericDescriptions[objClass] = ammoClipGenerics
}
