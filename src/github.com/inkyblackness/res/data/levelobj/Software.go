package levelobj

import (
	"github.com/inkyblackness/res/data/interpreters"
)

var baseSoftware = interpreters.New()

var multimediaFile = baseSoftware.
	With("ID", 1, 1).
	With("Type", 2, 1).As(interpreters.EnumValue(map[uint32]string{0: "E-Mail", 1: "Log", 2: "Data"}))

var cyberspaceProgram = baseSoftware.
	With("Version", 0, 1).As(interpreters.RangedValue(1, 9))

var funPack = baseSoftware.
	With("GameMask", 0, 1).As(interpreters.EnumValue(map[uint32]string{
	0x01: "Ping",
	0x02: "Eel Zapper",
	0x04: "Road",
	0x08: "Botbounce",
	0x10: "15",
	0x20: "TriopToe",
	0x80: "Wing 0"}))

var baseCyberspaceScenery = interpreters.New()

var scenerySoftware = baseCyberspaceScenery.
	With("Parameter", 0, 2).
	With("Subclass", 2, 4).As(interpreters.RangedValue(0, 7)).
	With("Type", 6, 4).As(interpreters.RangedValue(0, 16))

func initSoftware() interpreterRetriever {
	cyberspacePrograms := newInterpreterLeaf(cyberspaceProgram)

	realWorldTools := newInterpreterEntry(baseSoftware)
	realWorldTools.set(0, newInterpreterLeaf(funPack))

	multimediaSoftware := newInterpreterEntry(baseSoftware)
	multimediaSoftware.set(0, newInterpreterLeaf(multimediaFile)) // text
	multimediaSoftware.set(1, newInterpreterLeaf(multimediaFile)) // email

	class := newInterpreterEntry(baseSoftware)
	class.set(0, cyberspacePrograms) // aggressive programs
	class.set(1, cyberspacePrograms) // defensive programs
	class.set(2, cyberspacePrograms) // boost programs
	class.set(3, realWorldTools)
	class.set(4, multimediaSoftware)

	return class
}

func initCyberspaceScenery() interpreterRetriever {
	class := newInterpreterLeaf(scenerySoftware)

	return class
}
