package levelobj

import (
	"github.com/inkyblackness/res/data/interpreters"
)

var baseItem = interpreters.New()

var paperItem = baseItem.
	With("PaperId", 2, 1)

var briefcaseItem = baseItem.
	With("ObjectIndex1", 2, 2).As(interpreters.ObjectIndex()).
	With("ObjectIndex2", 4, 2).As(interpreters.ObjectIndex()).
	With("ObjectIndex3", 6, 2).As(interpreters.ObjectIndex()).
	With("ObjectIndex4", 8, 2).As(interpreters.ObjectIndex())

var corpseItem = baseItem.
	With("ObjectIndex1", 2, 2).As(interpreters.ObjectIndex()).
	With("ObjectIndex2", 4, 2).As(interpreters.ObjectIndex()).
	With("ObjectIndex3", 6, 2).As(interpreters.ObjectIndex()).
	With("ObjectIndex4", 8, 2).As(interpreters.ObjectIndex())

var severedHeadItem = baseItem.
	With("ImageIndex", 2, 1)

var accessCardItem = baseItem.
	With("AccessMask", 2, 4)

var securityIDModuleItem = baseItem.
	With("AccessMask", 2, 4)

var cyberInfoNodeItem = baseItem.
	With("TextIndex", 2, 1)

var cyberRestorative = baseItem.
	With("RestorationAmount", 2, 1)

var cyberDefenseMine = baseItem.
	With("DamageAmount", 2, 1)

var cyberBarricade = baseItem.
	With("Size", 2, 1).
	With("Height", 3, 1).
	With("Color", 6, 1)

func initItems() interpreterRetriever {

	junk := newInterpreterEntry(baseItem)
	junk.set(2, newInterpreterLeaf(paperItem))
	junk.set(7, newInterpreterLeaf(briefcaseItem))

	dead := newInterpreterEntry(baseItem)
	corpses := newInterpreterLeaf(corpseItem)
	severedHeads := newInterpreterLeaf(severedHeadItem)
	dead.set(0, corpses)
	dead.set(1, corpses)
	dead.set(2, corpses)
	dead.set(3, corpses)
	dead.set(4, corpses)
	dead.set(5, corpses)
	dead.set(6, corpses)
	dead.set(7, corpses)
	dead.set(13, severedHeads)
	dead.set(14, severedHeads)

	class := newInterpreterEntry(baseItem)
	class.set(0, junk)
	class.set(2, dead)
	class.set(4, newInterpreterLeaf(accessCardItem))

	return class
}

func initCyberspaceItems() interpreterRetriever {
	cyberspaceItems := newInterpreterEntry(baseItem)
	infoNodes := newInterpreterLeaf(cyberInfoNodeItem)
	cyberspaceItems.set(1, newInterpreterLeaf(cyberRestorative))
	cyberspaceItems.set(2, newInterpreterLeaf(cyberDefenseMine))
	cyberspaceItems.set(3, newInterpreterLeaf(securityIDModuleItem))
	cyberspaceItems.set(6, infoNodes)
	cyberspaceItems.set(8, infoNodes)
	cyberspaceItems.set(9, newInterpreterLeaf(cyberBarricade))

	class := newInterpreterEntry(baseItem)
	class.set(5, cyberspaceItems)

	return class
}
