package levelobj

import (
	"github.com/inkyblackness/res/data/interpreters"
)

var baseExplosive = interpreters.New().
	With("Unknown0000", 0, 2).
	With("State", 2, 2).
	With("TimerTime", 4, 2)

func initExplosives() interpreterRetriever {
	timedExplosives := newInterpreterEntry(baseExplosive)

	timedExplosives.set(2, newInterpreterLeaf(interpreters.New())) // Object explosion - not encountered

	class := newInterpreterEntry(baseExplosive)
	class.set(1, timedExplosives)

	return class
}
