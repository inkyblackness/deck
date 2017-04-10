package levelobj

import (
	"github.com/inkyblackness/res/data/interpreters"
)

var baseWeapon = interpreters.New()

var energyWeapon = baseWeapon.
	With("Charge", 0, 1).As(interpreters.RangedValue(0, 255)).
	With("Temperature", 1, 1).As(interpreters.RangedValue(0, 255))

var projectileWeapon = baseWeapon.
	With("AmmoType", 0, 1).As(interpreters.EnumValue(map[uint32]string{0: "Standard", 1: "Special"})).
	With("AmmoCount", 1, 1).As(interpreters.RangedValue(0, 255))

func initWeapons() interpreterRetriever {
	projectileWeapons := newInterpreterLeaf(projectileWeapon)
	energyWeapons := newInterpreterLeaf(energyWeapon)

	class := newInterpreterEntry(baseWeapon)
	class.set(0, projectileWeapons)
	class.set(1, projectileWeapons)
	class.set(2, projectileWeapons)

	class.set(4, energyWeapons)
	class.set(5, energyWeapons)

	return class
}
