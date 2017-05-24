package gameobj

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/data/interpreters"
)

var explosiveGenerics = interpreters.New().
	Refining("BasicWeapon", 0, 8, basicWeapon, interpreters.Always).
	With("BlastRange", 9, 1).
	With("BlastCoreRange", 10, 1).
	With("BlastDamage", 11, 1).
	With("ImpactForce", 12, 1).
	With("ExplosiveFlags", 13, 1)

var timedExplosives = interpreters.New().
	With("MinimumTime", 0, 1).
	With("MaximumTime", 1, 1).
	With("RandomFactor", 2, 1)

func initExplosives() {
	objClass := res.ObjectClass(3)

	genericDescriptions[objClass] = explosiveGenerics

	setSpecific(objClass, 1, timedExplosives)
}
