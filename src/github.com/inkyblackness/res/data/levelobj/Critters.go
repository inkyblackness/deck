package levelobj

import (
	"github.com/inkyblackness/res/data/interpreters"
)

var baseCritter = interpreters.New().
	With("StateTimeout", 0x0C, 2).
	With("PrimaryState", 0x15, 1).As(interpreters.EnumValue(map[uint32]string{
	0: "docile",
	1: "cautious",
	2: "hostile",
	3: "cautious (?)",
	4: "attacking",
	5: "sleeping",
	6: "tranquilized",
	7: "confused"})).
	With("SecondaryState", 0x16, 1).As(interpreters.EnumValue(map[uint32]string{
	0: "docile",
	1: "cautious",
	2: "hostile",
	3: "cautious (?)",
	4: "attacking",
	5: "sleeping",
	6: "tranquilized",
	7: "confused"})).
	With("LootObjectIndex1", 0x20, 2).As(interpreters.ObjectIndex()).
	With("LootObjectIndex2", 0x22, 2).As(interpreters.ObjectIndex())

func initCritters() interpreterRetriever {
	class := newInterpreterEntry(baseCritter)

	return class
}
