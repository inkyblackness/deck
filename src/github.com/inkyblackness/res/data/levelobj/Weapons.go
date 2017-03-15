package levelobj

import (
	"github.com/inkyblackness/res/data/interpreters"
)

var baseWeapon = interpreters.New()

var energyWeapon = baseWeapon.
	With("Charge", 0, 1).
	With("Temperature", 1, 1)

var projectileWeapon = baseWeapon.
	With("AmmoType", 0, 1).
	With("AmmoCount", 1, 1)

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
