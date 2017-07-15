package levelobj

import (
	"github.com/inkyblackness/res/data/interpreters"
)

var realWorldEntries *interpreterEntry
var cyberspaceEntries *interpreterEntry

var realWorldExtras *interpreterEntry
var cyberspaceExtras *interpreterEntry

var extraIced = interpreters.New().
	With("ICE-presence", 1, 1).
	With("ICE-level", 3, 1)

func init() {

	projectiles := newInterpreterEntry(baseProjectile)
	software := initSoftware()
	animations := newInterpreterEntry(baseAnimation)
	markers := initMarkers()
	critters := initCritters()

	realWorldEntries = newInterpreterEntry(interpreters.New())
	realWorldEntries.set(0, initWeapons())
	realWorldEntries.set(1, newInterpreterEntry(interpreters.New())) // have no data
	realWorldEntries.set(2, projectiles)
	realWorldEntries.set(3, initExplosives())
	realWorldEntries.set(4, newInterpreterEntry(interpreters.New())) // have no data
	realWorldEntries.set(5, newInterpreterEntry(baseHardware))
	realWorldEntries.set(6, software)
	realWorldEntries.set(7, initScenery())
	realWorldEntries.set(8, initItems())
	realWorldEntries.set(9, initPanels())
	realWorldEntries.set(10, initBarriers())
	realWorldEntries.set(11, animations)
	realWorldEntries.set(12, markers)
	realWorldEntries.set(13, initContainers())
	realWorldEntries.set(14, critters)

	cyberspaceEntries = newInterpreterEntry(interpreters.New())
	cyberspaceEntries.set(6, software)
	cyberspaceEntries.set(7, initCyberspaceScenery())
	cyberspaceEntries.set(8, initCyberspaceItems())
	cyberspaceEntries.set(9, initCyberspacePanels())
	cyberspaceEntries.set(12, markers)
	cyberspaceEntries.set(14, critters)

	realWorldExtras = newInterpreterEntry(interpreters.New())
	cyberspaceExtras = newInterpreterEntry(interpreters.New())
	cyberspaceExtras.set(6, newInterpreterLeaf(extraIced))
	cyberspaceExtras.set(7, newInterpreterLeaf(extraIced))
	cyberspaceExtras.set(8, newInterpreterLeaf(extraIced))
	cyberspaceExtras.set(9, newInterpreterLeaf(extraIced))
}
