package levelobj

import (
	"github.com/inkyblackness/res/data/interpreters"
)

var baseSoftware = interpreters.New()

var multimediaFile = baseSoftware.
	With("ID", 1, 1).
	With("Type", 2, 1)

var cyberspaceProgram = baseSoftware.
	With("Version", 0, 1)

var funPack = baseSoftware.
	With("GameMask", 0, 1)

var baseCyberspaceScenery = interpreters.New()

var scenerySoftware = baseCyberspaceScenery.
	With("Parameter", 0, 2).
	With("Subclass", 2, 4).
	With("Type", 6, 4)

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
